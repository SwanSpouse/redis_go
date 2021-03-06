package server

import (
	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/handlers"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/rdb"
	"github.com/SwanSpouse/redis_go/util"
)

var (
	_ rdb.IDecoder = (*Server)(nil)
)

func (srv *Server) StartRDB() {
	loggers.Info("rdb redis rdb started")
}

func (srv *Server) StartDatabase(n int) {
	loggers.Info("current database NO: %d\n", n)
	srv.dbIndex = n
	srv.FakeClient.SetDatabase(srv.Databases[n])
}

func (srv *Server) Aux(key, value []byte) {
	loggers.Info("rdb process aux key:%s value:%s", key, value)
}

func (srv *Server) ResizeDatabase(dbSize, expiresSize uint32) {
	loggers.Info("rdb process resize database dbSize:%d expiresSize:%d", dbSize, expiresSize)

}

func (srv *Server) Set(key, value []byte, expiry int64) {
	loggers.Info("rdb process set key:%s value:%s", key, value)
	srv.FakeClient.Argv = []string{handlers.RedisStringCommandSet, string(key), string(value)}
	srv.commandTable[handlers.RedisStringCommandSet].Proc(srv.FakeClient)
}

func (srv *Server) StartHash(key []byte, length, expiry int64) {
	loggers.Info("rdb process start hash key:%s length:%d", key, length)
}

func (srv *Server) Hset(key, field, value []byte) {
	loggers.Info("rdb process HSet key:%s field:%s, value:%s", key, field, value)
	srv.FakeClient.Argv = []string{handlers.RedisHashCommandHSet, string(key), string(field), string(value)}
	srv.commandTable[handlers.RedisHashCommandHSet].Proc(srv.FakeClient)
}

func (srv *Server) EndHash(key []byte) {
	loggers.Info("rdb process end hash key:%s", key)
}

func (srv *Server) StartSet(key []byte, cardinality, expiry int64) {
	loggers.Info("rdb process start set key:%s", key)
}

func (srv *Server) Sadd(key, member []byte) {
	loggers.Info("rdb process SAdd key:%s, member:%s", key, member)
	// TODO lmj
}

func (srv *Server) EndSet(key []byte) {
	loggers.Info("rdb process start set key:%s", key)
}

func (srv *Server) StartList(key []byte, length, expiry int64) {
	loggers.Info("rdb process start list key%s", key)
}

func (srv *Server) Rpush(key, value []byte) {
	loggers.Info("rdb process RPush key:%s value%s", key, value)
	srv.FakeClient.Argv = []string{handlers.RedisListCommandRPush, string(key), string(value)}
	srv.commandTable[handlers.RedisListCommandRPush].Proc(srv.FakeClient)
}

func (srv *Server) EndList(key []byte) {
	loggers.Info("rdb process end list key:%s", key)
}

func (srv *Server) StartZSet(key []byte, cardinality, expiry int64) {
	loggers.Info("rdb process start ZSet key:%s", key)

}

func (srv *Server) Zadd(key []byte, score float64, member []byte) {
	loggers.Info("rdb process ZAdd key:%s", key)
	// TODO lmj
}

func (srv *Server) EndZSet(key []byte) {
	loggers.Info("rdb process End ZSet key:%s", key)

}

func (srv *Server) EndDatabase(n int) {
	loggers.Info("rdb process End Database db:%d", n)

}

func (srv *Server) EndRDB() {
	loggers.Info("rdb process End RDB")
}

func (srv *Server) rdbLoad() {
	if !util.FileExists(srv.Config.RdbFilename) {
		loggers.Info("redis rdb file not exits")
		return
	}
	srv.FakeClient = client.NewFakeClient()
	decoder, err := rdb.NewDecoder(srv.Config.RdbFilename, srv)
	if err != nil {
		loggers.Errorf("rdb new encoder error %+v", err)
		return
	}
	decoder.Decode()
	if err != nil {
		loggers.Errorf("rdb decode error %+v", err)
	}
}
