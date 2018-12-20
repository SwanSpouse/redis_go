package server

import (
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/conf"
	"github.com/SwanSpouse/redis_go/database"
	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/tcp"
	"github.com/SwanSpouse/redis_go/util"
)

const (
	RedisServerStatusNormal             = 0
	RedisServerStatusRdbSaveInProcess   = 1
	RedisServerStatusRdbBgSaveInProcess = 2
)

// Redis server
type Server struct {
	TcpListener      net.Listener
	clientIDSequence int64 // client auto increasing sequence id
	Config           *conf.ServerConfig
	Databases        []*database.Database     /* database*/
	dbIndex          int                      // rdb process current db
	clients          map[int64]*client.Client // clientID -> client
	FakeClient       *client.Client           // used in rdb and aof
	password         string                   /* Pass for AUTH command, or NULL */
	commandTable     map[string]*client.Command
	mu               sync.RWMutex
	Status           atomic.Value
	Dirty            int64
	rdbLastSave      time.Time
	aofSelectDBId    int
	aofLock          sync.Mutex // aof lock
	aofBuf           []byte     // append only file buffer
	aofLastSave      time.Time  // aof last save time
	TimeEventLoop    *EventLoop // redis time event
	WaitGroup        util.WaitGroupWrapper
	ExitChan         chan int
}

func NewServer(config *conf.ServerConfig) *Server {
	server := &Server{
		Config:        config,
		commandTable:  make(map[string]*client.Command),
		clients:       make(map[int64]*client.Client),
		TimeEventLoop: NewEventLoop(),
		ExitChan:      make(chan int),
	}
	// init general parameters
	server.initServer()

	// init Reader & Writer Sync Pool
	server.initIOPool()

	// init databases
	server.initDB()

	// init commandTable table
	server.populateCommandTable()

	// 在这里把 serverCron 添加到timeEvent里面
	server.TimeEventLoop.NewTimeEvent(100, 0, true, server.ServerCron)

	// load data
	server.loadDataFromDisk()

	loggers.Debug("redis server: %+v", server)
	return server
}

// 判断server 此时是否可以对外提供服务
func (srv *Server) isServiceAvailable() bool {
	return srv.Status.Load() == RedisServerStatusNormal || srv.Status.Load() == RedisServerStatusRdbBgSaveInProcess
}

// 处理来自客户端的请求
func (srv *Server) IOLoop(conn net.Conn) {
	loggers.Info("TCP: new client(%s)", conn.RemoteAddr())
	// TODO lmj srv.clients 有个数限制

	c := client.NewClient(atomic.AddInt64(&srv.clientIDSequence, 1), conn, srv.getDefaultDB())
	srv.addClient(c)

	var err error
	// handle client command
	for {
		// read command from client
		if err = c.ProcessInputBuffer(); err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				loggers.Errorf("server read command error %+v", err)
				c.ResponseReError(err)
				continue
			}
		}
		if !srv.isServiceAvailable() {
			c.ResponseReError(re.ErrRedisRdbSaveInProcess)
			continue
		}
		/**
		首先判断是否在command table中,
			如果不在command table中,则返回command not found
			如果在command table中，则获取到相应的command handler来进行处理。
		*/
		command, ok := srv.commandTable[strings.ToUpper(c.Argv[0])]
		if !ok || command == nil {
			loggers.Errorf(string(re.ErrUnknownCommand), c.Argv[0])
			c.ResponseReError(re.ErrUnknownCommand, c.Argv[0])
			continue
		}
		c.LastCmd = c.Cmd
		c.Cmd = command
		/**
		在这里对command的参数个数等进行检查
			1. 如果Arity > 0, 要求参数个数必须严格等于Arity
			2. 如果Arity < 0, 要求参数个数至少为|Arity|
		*/
		if (command.Arity > 0 && c.Argc != command.Arity) || (c.Argc < -command.Arity) {
			loggers.Errorf("wrong number of args %+v", command)
			c.ResponseReError(re.ErrWrongNumberOfArgs, c.Argv[0])
			continue
		}

		// TODO 检查用户是否验证过身份
		// TODO 集群模式等在这里进行一些操作
		// TODO 判断是否是事务相关命令
		// TODO 判断命令造成了多少个dirty, 执行时间等一些统计信息
		// 在这里对client端发送过来的命令进行处理
		command.Handler.Process(c)

		// 在rdb save结束之后，重新统计dirty数量并记录本次rdb结束的时间
		if c.Cmd.GetName() == RedisServerCommandSave {
			srv.Dirty = 0
			srv.rdbLastSave = time.Now()
		} else {
			srv.Dirty += c.Dirty
		}

		// 在这里判断命令是否要发送到aof_buf或者Aof文件
		if srv.Config.AofState == conf.RedisAofOn && c.Cmd.Flags&client.RedisCmdWrite > 0 && c.Dirty != 0 {
			loggers.Debug("Client exec a write cmd or make db dirty")
			srv.propagate(c)
			// 现在默认将每个写命令都刷写到aof文件中
			srv.flushAppendOnlyFile()
			c.Dirty = 0
		}
	}
	loggers.Info("client %d-%s exiting ioLoop", c.ID(), c.RemoteAddr())
	if err != nil {
		loggers.Errorf("client %d %s", c.ID(), err)
	}
	// remove client form server
	srv.removeClient(c)
}

// 启动redis server 并开始监听TCP连接
func (srv *Server) TCPServe() {
	loggers.Info("TCP: listening on %s", srv.TcpListener.Addr())

	// loop for accepting tcp connection from redis client
	for {
		// 接收来自客户端的请求
		clientConn, err := srv.TcpListener.Accept()
		if err != nil {
			// ignore temporary errors
			if netError, ok := err.(net.Error); ok && netError.Temporary() {
				loggers.Warn("temporary Accept() failure %s", err)
				runtime.Gosched()
				continue
			}
			break
		}
		go srv.IOLoop(clientConn)
	}
	loggers.Info("TCP: closing %s", srv.TcpListener.Addr())
}

func (srv *Server) addClient(c *client.Client) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	srv.clients[c.ID()] = c
}

func (srv *Server) removeClient(c *client.Client) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	c.Close()
	client.ReturnClient(c)
	delete(srv.clients, c.ID())
}

func (srv *Server) initServer() {
	srv.Status.Store(RedisServerStatusNormal)
	srv.aofSelectDBId = -1
	if srv.Config.AofState == conf.RedisAofOn {
		srv.aofBuf = make([]byte, 0)
	}
	loggers.Level = srv.Config.LogLevel
}

func (srv *Server) initDB() {
	// add default database
	srv.Databases = make([]*database.Database, srv.Config.DBNum)
	for i := 0; i < srv.Config.DBNum; i++ {
		srv.Databases[i] = database.NewDatabase(i)
	}
}

func (srv *Server) initIOPool() {
	tcp.InitBufIoReaderPool(srv.Config.ReaderPoolSize)
	tcp.InitBufIoWriterPool(srv.Config.WriterPoolSize)
	loggers.Debug("Successful init reader and writer pool. ReaderPoolSize:%d, WriterPoolSize:%d", srv.Config.ReaderPoolSize, srv.Config.WriterPoolSize)
}

func (srv *Server) processTimeEvents() {
	loggers.Debug("start to processTimeEvents")
	workTicker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-workTicker.C:
			srv.TimeEventLoop.lock.Lock()

			currentTime := util.GetCurrentMillisecond()
			for timeEventId, timeEvent := range srv.TimeEventLoop.events {
				// 超过时间间隔，需要被执行
				if currentTime-srv.TimeEventLoop.lastTime >= timeEvent.Interval {
					timeEvent.Proc()
					// 执行一次的任务在执行过后进行删除
					if !timeEvent.runInCircle {
						delete(srv.TimeEventLoop.events, timeEventId)
					}
				}
			}
			srv.TimeEventLoop.lastTime = currentTime
			srv.TimeEventLoop.lock.Unlock()
		case <-srv.ExitChan:
			goto exit
		}
	}
exit:
	loggers.Info("PROCESS TIME EVENTS: closing")
	workTicker.Stop()
}

func (srv *Server) getDefaultDB() *database.Database {
	if srv.Databases == nil {
		srv.initDB()
	}
	return srv.Databases[0]
}

func (srv *Server) loadDataFromDisk() {
	startTime := time.Now()
	if srv.Config.AofState == conf.RedisAofOn {
		loggers.Info("redis aof start to load data from disk at %s", startTime.Format("20060102 15:04:05"))
		srv.loadAppendOnlyFile()
	} else {
		loggers.Info("redis rdb start to load data from disk at %s", startTime.Format("20060102 15:04:05"))
		srv.rdbLoad()
	}
}

/**
ServerCron主要负责一下工作:
	更新服务器的各类统计信息，比如时间、内存占用、数据库占用情况等。
	清理数据库中的过期键值对。
	对不合理的数据库进行大小调整。
	关闭和清理连接失效的客户端。
	尝试进行 AOF 或 RDB 持久化操作。
	如果服务器是主节点的话，对附属节点进行定期同步。
	如果处于集群模式的话，对集群进行定期同步和连接测试。
*/
func (srv *Server) ServerCron() {
}
