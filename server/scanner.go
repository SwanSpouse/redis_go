package server

import (
	"io"
	"redis_go/client"
	"redis_go/conf"
	re "redis_go/error"
	"redis_go/loggers"
	"strings"
	"sync/atomic"
	"time"
)

func (srv *Server) scanClients() {
	for {
		srv.mu.Lock()
		for i := 0; i < len(srv.clients); i++ {
			if srv.clients[i] == nil || srv.clients[i].Closed {
				srv.clients[i] = nil
				continue
			}
			//if srv.clients[i].Status == client.RedisClientStatusInProcess && srv.clients[i].IsExecTimeout() {
			//	log.Info("redis client %d execute timeout at %+v", i, srv.clients[i].GetExecTimeoutAt())
			//	srv.clients[i].Close()
			//	continue
			//}
			//if srv.clients[i].Status == client.RedisClientStatusIdle && srv.clients[i].IsIdleTimeout() {
			//	log.Info("redis client %d idle timeout at %+v", i, srv.clients[i].GetIdleTimeoutAt())
			//	srv.clients[i].Close()
			//	continue
			//}
			// 如果Client有待处理命令，处理对应的Client
			if srv.clients[i].Status == client.RedisClientStatusIdle {
				atomic.StoreUint32(&srv.clients[i].Status, client.RedisClientStatusInProcess)
				go srv.handlerCommand(srv.clients[i])
			}
		}
		srv.mu.Unlock()
	}
}

func (srv *Server) handlerCommand(c *client.Client) {
	if c == nil || c.Closed {
		return
	}
	c.Locker.Lock()
	defer c.Locker.Unlock()

	defer atomic.StoreUint32(&c.Status, client.RedisClientStatusIdle)
	//c.SetIdleTimeout(5 * time.Hour)
	//c.SetExecTimeout(5 * time.Second)

	loggers.Info("handler command from client:%s", c.RemoteAddr())
	// ReadCmd这里会阻塞知道有数据或者客户端断开连接
	if err := c.ProcessInputBuffer(); err != nil && err == io.EOF {
		c.Close()
		return
	} else if err != nil {
		loggers.Errorf("server read command error %+v", err)
		c.ResponseReError(err)
		return
	}

	// 如果服务器正在进行阻塞操作，不接受客户端发过来的请求
	if !srv.isInService() {
		c.ResponseReError(re.ErrRedisRdbSaveInProcess)
		return
	}
	/**
	首先判断是否在command table中,
		如果不在command table中,则返回command not found
		如果在command table中，则获取到相应的command handler来进行处理。
	*/
	if command, ok := srv.commandTable[strings.ToUpper(c.Argv[0])]; !ok || command == nil {
		loggers.Errorf(string(re.ErrUnknownCommand), c.Argv[0])
		c.ResponseReError(re.ErrUnknownCommand, c.Argv[0])
	} else {
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
			return
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

			// 将srv aof_buf中的数据同步到文件中
			//srv.flushAppendOnlyFile()
			c.Dirty = 0
		}
	}
}
