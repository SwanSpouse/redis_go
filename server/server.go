package server

import (
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/database"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/tcp"
	"sync"
)

// Redis server
type Server struct {
	Config       *conf.ServerConfig
	Databases    []*database.Database /* database*/
	clients      []*client.Client
	password     string /* Pass for AUTH command, or NULL */
	commandTable map[string]*client.Command
	mu           sync.RWMutex
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = conf.InitServerConfig()
	}
	server := &Server{
		Config:       config,
		commandTable: make(map[string]*client.Command),
	}
	loggers.Info("redis server config: %+v", config)
	// init general parameters
	server.initServer()
	// init Reader & Writer Sync Pool
	server.initIOPool()
	// init databases
	server.initDB()
	// init commandTable table
	server.populateCommandTable()
	// init time events
	go server.initTimeEvents()
	loggers.Info("redis server: %+v", server)
	return server
}

func (srv *Server) Serve(lis net.Listener) error {
	// start to scan clients
	go srv.scanClients()
	// loop for accept tcp client
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		c := client.NewClient(cn, srv.getDefaultDB())
		srv.addClientToServer(c)
		loggers.Info("new client %d come in ! from %+v and has been added in server's client list.", c.ID(), cn.RemoteAddr().String())
	}
}

func (srv *Server) addClientToServer(c *client.Client) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// add client to server's clients list
	for i := 0; i < len(srv.clients); i++ {
		if srv.clients[i] == nil || srv.clients[i].Closed {
			srv.clients[i] = c
			return
		}
	}
	// TODO lmj srv.clients 有个数限制
	srv.clients = append(srv.clients, c)
}

func (srv *Server) initServer() {
	loggers.Level = srv.Config.LogLevel
}

func (srv *Server) initDB() {
	// add default database
	srv.Databases = make([]*database.Database, srv.Config.DBNum)
	for i := 0; i < srv.Config.DBNum; i++ {
		srv.Databases[i] = database.NewDatabase()
	}
}

func (srv *Server) initIOPool() {
	for i := 0; i < srv.Config.ReaderPoolNum; i++ {
		tcp.ReaderPool.Put(tcp.NewBufIoReaderWithoutConn())
	}
	for i := 0; i < srv.Config.WriterPoolNum; i++ {
		tcp.WriterPool.Put(tcp.NewBufIoWriterWithoutConn())
	}
	loggers.Debug("Successful init reader and writer pool. ReaderPoolSize:%d, WriterPoolSize:%d", srv.Config.ReaderPoolNum, srv.Config.WriterPoolNum)
}

func (srv *Server) initTimeEvents() {
	//ticker := time.NewTicker(time.Second)
	//for _ = range ticker.C {
	//	log.Info("TICKER INFO client list length %d", len(srv.clients))
	//	for i, item := range srv.clients {
	//		log.Info("TICKER current client list index %d, info:%+v", i, item)
	//	}
	//}
}

func (srv *Server) getDefaultDB() *database.Database {
	if srv.Databases == nil {
		srv.initDB()
	}
	return srv.Databases[0]
}

/*
 * This is the meaning of the flags:
 *
 * w: write command (may modify the key space).
 *    写入命令，可能会修改 key space
 *
 * r: read command  (will never modify the key space).
 *    读命令，不修改 key space
 *
 * m: may increase memory usage once called. Don't allow if out of memory.
 *    可能会占用大量内存的命令，调用时对内存占用进行检查
 *
 * a: admin command, like SAVE or SHUTDOWN.
 *    管理员使用的命令
 *
 * p: Pub/Sub related command.
 *    发送/订阅相关的命令
 *
 * f: force replication of this command, regarless of server.dirty.
 *    强制同步这个命令，无视 server.dirty
 *
 * s: command not allowed in scripts.
 *    不允许在脚本中使用的命令
 *
 * R: random command. Command is not deterministic, that is, the same command
 *    with the same arguments, with the same key space, may have different
 *    results. For instance SPOP and RANDOMKEY are two random commands.
 *    随机命令，对于同样数据集的同一个命令调用，得出的结果可能是不相同的。
 *
 * S: Sort command output array if called from script, so that the output
 *    is deterministic.
 *    如果命令在脚本中执行，那么对输出进行排序，从而让输出变得确定起来。
 *
 * l: Allow command while loading the database.
 *    允许在载入数据库时执行的命令
 *
 * t: Allow command while a slave has stale data but is not allowed to
 *    server this data. Normally no command is accepted in this condition
 *    but just a few.
 *    允许在附属节点包含过期数据时执行的命令
 *
 * M: Do not automatically propagate the command on MONITOR.
 *    不要自动将此命令发送到 MONITOR
 */
func (srv *Server) populateCommandTable() {
	connectionHandler := new(handlers.ConnectionHandler)
	stringHandler := new(handlers.StringHandler)
	keyHandler := new(handlers.KeyHandler)
	listHandler := new(handlers.ListHandler)
	hashHandler := new(handlers.HashHandler)

	// connection command
	srv.commandTable[handlers.RedisConnectionCommandPing] = client.NewCommand(handlers.RedisConnectionCommandPing, 1, "r", connectionHandler)
	srv.commandTable[handlers.RedisConnectionCommandAuth] = client.NewCommand(handlers.RedisConnectionCommandAuth, 2, "rs", connectionHandler)

	// key command
	srv.commandTable[handlers.RedisKeyCommandDel] = client.NewCommand(handlers.RedisKeyCommandDel, -2, "w", keyHandler)
	srv.commandTable[handlers.RedisKeyCommandObject] = client.NewCommand(handlers.RedisKeyCommandObject, -2, "r", keyHandler)
	srv.commandTable[handlers.RedisKeyCommandType] = client.NewCommand(handlers.RedisKeyCommandType, 2, "r", keyHandler)
	srv.commandTable[handlers.RedisKeyCommandExists] = client.NewCommand(handlers.RedisKeyCommandExists, 2, "r", keyHandler)

	// string command
	srv.commandTable[handlers.RedisStringCommandAppend] = client.NewCommand(handlers.RedisStringCommandAppend, 1, "r", stringHandler)
	srv.commandTable[handlers.RedisStringCommandSet] = client.NewCommand(handlers.RedisStringCommandSet, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandMSet] = client.NewCommand(handlers.RedisStringCommandMSet, -3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandMSetNx] = client.NewCommand(handlers.RedisStringCommandMSetNx, -3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandSetNx] = client.NewCommand(handlers.RedisStringCommandSetNx, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandGet] = client.NewCommand(handlers.RedisStringCommandGet, 2, "r", stringHandler)
	srv.commandTable[handlers.RedisStringCommandMGet] = client.NewCommand(handlers.RedisStringCommandMGet, -2, "r", stringHandler)
	srv.commandTable[handlers.RedisStringCommandGetSet] = client.NewCommand(handlers.RedisStringCommandGetSet, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandIncr] = client.NewCommand(handlers.RedisStringCommandIncr, 2, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandIncrBy] = client.NewCommand(handlers.RedisStringCommandIncrBy, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandIncrByFloat] = client.NewCommand(handlers.RedisStringCommandIncrByFloat, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandDecr] = client.NewCommand(handlers.RedisStringCommandDecr, 2, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandDecrBy] = client.NewCommand(handlers.RedisStringCommandDecrBy, 3, "wm", stringHandler)
	srv.commandTable[handlers.RedisStringCommandStrLen] = client.NewCommand(handlers.RedisStringCommandStrLen, 2, "r", stringHandler)

	// list command
	srv.commandTable[handlers.RedisListCommandLIndex] = client.NewCommand(handlers.RedisListCommandLIndex, 3, "r", listHandler)
	srv.commandTable[handlers.RedisListCommandLInsert] = client.NewCommand(handlers.RedisListCommandLInsert, 5, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandLLen] = client.NewCommand(handlers.RedisListCommandLLen, 2, "r", listHandler)
	srv.commandTable[handlers.RedisListCommandLPop] = client.NewCommand(handlers.RedisListCommandLPop, 2, "w", listHandler)
	srv.commandTable[handlers.RedisListCommandLPush] = client.NewCommand(handlers.RedisListCommandLPush, -3, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandLPushX] = client.NewCommand(handlers.RedisListCommandLPushX, 3, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandLRange] = client.NewCommand(handlers.RedisListCommandLRange, 4, "r", listHandler)
	srv.commandTable[handlers.RedisListCommandLRem] = client.NewCommand(handlers.RedisListCommandLRem, 4, "w", listHandler)
	srv.commandTable[handlers.RedisListCommandLSet] = client.NewCommand(handlers.RedisListCommandLSet, 4, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandLTrim] = client.NewCommand(handlers.RedisListCommandLTrim, 4, "w", listHandler)
	srv.commandTable[handlers.RedisListCommandRPop] = client.NewCommand(handlers.RedisListCommandRPop, 2, "w", listHandler)
	srv.commandTable[handlers.RedisListCommandRPopLPush] = client.NewCommand(handlers.RedisListCommandRPopLPush, 3, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandRPush] = client.NewCommand(handlers.RedisListCommandRPush, -3, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandRpushX] = client.NewCommand(handlers.RedisListCommandRpushX, 3, "wm", listHandler)
	srv.commandTable[handlers.RedisListCommandLDebug] = client.NewCommand(handlers.RedisListCommandLDebug, 2, "r", listHandler)

	// hash command
	srv.commandTable[handlers.RedisHashCommandHDel] = client.NewCommand(handlers.RedisHashCommandHDel, -3, "w", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHExists] = client.NewCommand(handlers.RedisHashCommandHExists, 3, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHGet] = client.NewCommand(handlers.RedisHashCommandHGet, 3, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHGetAll] = client.NewCommand(handlers.RedisHashCommandHGetAll, 2, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHIncrBy] = client.NewCommand(handlers.RedisHashCommandHIncrBy, 4, "wm", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHIncrByFloat] = client.NewCommand(handlers.RedisHashCommandHIncrByFloat, 4, "wm", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHKeys] = client.NewCommand(handlers.RedisHashCommandHKeys, 2, "rS", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHLen] = client.NewCommand(handlers.RedisHashCommandHLen, 2, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHMGet] = client.NewCommand(handlers.RedisHashCommandHMGet, -3, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHMSet] = client.NewCommand(handlers.RedisHashCommandHMSet, -4, "wm", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHSet] = client.NewCommand(handlers.RedisHashCommandHSet, 4, "wm", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHSetNX] = client.NewCommand(handlers.RedisHashCommandHSetNX, 4, "wm", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHVals] = client.NewCommand(handlers.RedisHashCommandHVals, 2, "rS", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHScan] = client.NewCommand(handlers.RedisHashCommandHScan, 3, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHStrLen] = client.NewCommand(handlers.RedisHashCommandHStrLen, 3, "r", hashHandler)
	srv.commandTable[handlers.RedisHashCommandHDebug] = client.NewCommand(handlers.RedisHashCommandHDebug, 2, "r", hashHandler)
}
