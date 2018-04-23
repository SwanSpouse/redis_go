package server

import (
	"io"
	"redis_go/client"
	re "redis_go/error"
	"redis_go/loggers"
	"sync/atomic"
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
			go srv.handlerCommand(srv.clients[i])
		}
		srv.mu.Unlock()
	}
}

func (srv *Server) handlerCommand(c *client.Client) {
	if c == nil || c.Closed || c.Status == client.RedisClientStatusInProcess {
		return
	}
	c.Locker.Lock()
	defer c.Locker.Unlock()

	atomic.StoreUint32(&c.Status, client.RedisClientStatusInProcess)
	defer atomic.StoreUint32(&c.Status, client.RedisClientStatusIdle)

	//c.SetIdleTimeout(5 * time.Hour)
	//c.SetExecTimeout(5 * time.Second)
	if cmd, err := c.ReadCmd(); err == io.EOF {
		c.Close()
		return
	} else if err != nil || cmd == nil {
		return
	}
	/**
	首先判断是否在command table中,
		如果不在command table中,则返回command not found
		如果在command table中，则获取到相应的command handler来进行处理。
	*/
	if handler, ok := srv.commands[c.Cmd.GetName()]; ok {
		/* 在这里对client端发送过来的命令进行处理 */
		handler.Process(c)
	} else {
		loggers.Errorf(string(re.ErrUnknownCommand), c.Cmd.GetOriginName())
		c.ResponseReError(re.ErrUnknownCommand, c.Cmd.GetOriginName())
	}
}