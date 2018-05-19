package server

import (
	"fmt"
	"io"
	"redis_go/aof"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/handlers"
	"redis_go/loggers"
	"strings"
	"time"
)

func appendStrToByteArr(input []byte, args ...string) []byte {
	for _, arg := range args {
		input = append(input, []byte(arg)...)
	}
	return input
}

func catAppendOnlyGenericCommand(buf []byte, argc int, argv []string) []byte {
	// 先处理参数个数
	buf = appendStrToByteArr(buf, "*")
	buf = appendStrToByteArr(buf, fmt.Sprintf("%d\r\n", argc))

	for _, arg := range argv {
		buf = appendStrToByteArr(buf, fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return buf
}

func (srv *Server) propagate(c *client.Client) {
	loggers.Debug("propagate cmd to server aof buf")

	outBuf := make([]byte, 0)
	dbId := c.SelectedDatabase().GetID()
	if dbId != srv.aofSelectDBId {
		dbStr := fmt.Sprintf("%d", dbId)
		outBuf = appendStrToByteArr(outBuf, fmt.Sprintf("*2\r\n$6\r\nSELECT\r\n$%d\r\n%d\r\n", len(dbStr), dbId))
		srv.aofSelectDBId = c.SelectedDatabase().GetID()
	}

	if c.Cmd.GetName() == handlers.RedisKeyCommandPExpire {
		// TODO lmj
	} else if c.Cmd.GetName() == handlers.RedisStringCommandSetEX || c.Cmd.GetName() == handlers.RedisStringCommandPSetEx {
		// TODO lmj
	} else {
		outBuf = catAppendOnlyGenericCommand(outBuf, c.Argc, c.Argv)
	}

	// 将命令写入srv的 aof_buf，下次同步到aof文件的时候这些数据就会被刷新到文件中。
	srv.aofBuf = append(srv.aofBuf, outBuf...)
	loggers.Debug("current aof bug:%s", string(srv.aofBuf))
}

func (srv *Server) flushAppendOnlyFile() {
	if len(srv.aofBuf) == 0 {
		return
	}

	if srv.Config.AofFSync == conf.RedisAofFSyncEverySec {
		// TODO 这里有策略可以进行延迟写
		loggers.Info("Hi~ You have a todo here. ")
		return
	}
	loggers.Debug("start to flush aof file")

	encoder, err := aof.NewEncoder(srv.Config.AofFilename)
	if err != nil || encoder == nil {
		loggers.Errorf("new aof encoder error:%+v", err)
	}
	if n, err := encoder.Write(srv.aofBuf); err != nil {
		loggers.Errorf("flush aof data to file error:%+v", err)
	} else if n != len(srv.aofBuf) {
		loggers.Errorf("number:%d of written data is not equal with aof buf length:%d", n, len(srv.aofBuf))
	} else {
		srv.aofLastSave = time.Now()
		srv.aofBuf = make([]byte, 0)
		loggers.Debug("flush aof file end")
	}
}

func (srv *Server) loadAppendOnlyFile() {
	loggers.Debug("start to load append only file")
	decoder := aof.NewDecoder(srv.Config.AofFilename)
	if decoder == nil {
		loggers.Info("aof file:%s not exists", srv.Config.AofFilename)
		return
	}
	// 创建伪终端来发送命令
	srv.FakeClient = client.NewFakeClient()
	srv.FakeClient.SetDatabase(srv.Databases[0])
	for true {
		out, err := decoder.DecodeAppendOnlyFile()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			loggers.Errorf("load append only file error:%+v", err)
			return
		}
		cmd, ok := srv.commandTable[strings.ToUpper(out.Argv[0])]
		if !ok {
			loggers.Errorf("Unknown command '%s' reading the append only file", out.Argv[0])
			return
		}
		if (cmd.Arity > 0 && out.Argc != cmd.Arity) || (out.Argc < -cmd.Arity) {
			loggers.Errorf("wrong number of args %+v", out.Argv[0])
			return
		}
		loggers.Debug("current cmd we receive in aof argv:%+v", out.Argv)
		srv.FakeClient.LastCmd = srv.FakeClient.Cmd
		srv.FakeClient.Cmd = cmd
		srv.FakeClient.Argc = out.Argc
		srv.FakeClient.Argv = out.Argv
		// process command
		cmd.Handler.Process(srv.FakeClient)
		srv.FakeClient.Argc = 0
		srv.FakeClient.Argv = nil
	}
	loggers.Debug("load append only file end")
}
