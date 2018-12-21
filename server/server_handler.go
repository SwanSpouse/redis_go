package server

import (
	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/database"
	"github.com/SwanSpouse/redis_go/encodings"
	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/rdb"
)

/*
	这里这么写虽然和handlers里面那些逻辑有些不符。但Server 本身就是一个能够处理命令的handler
    这样想想也是合理的。
*/

const (
	RedisServerCommandBGSRewriteAof = "BGREWRITEAOF"
	RedisServerCommandBGSave        = "BGSAVE"
	RedisServerCommandClient        = "CLIENT"
	RedisServerCommandConfig        = "CONFIG"
	RedisServerCommandDBSize        = "DBSIZE"
	RedisServerCommandDebug         = "DEBUG"
	RedisServerCommandFlushAll      = "FLUSHALL"
	RedisServerCommandFlushDB       = "FLUSHDB"
	RedisServerCommandInfo          = "INFO"
	RedisServerCommandLastSave      = "LASTSAVE"
	RedisServerCommandMonitor       = "MONITOR"
	RedisServerCommandPSync         = "PSYNC"
	RedisServerCommandSave          = "SAVE"
	RedisServerCommandShutDown      = "SHUTDOWN"
	RedisServerCommandSlaveOf       = "SLAVEOF"
	RedisServerCommandSlowLog       = "SLOWLOG"
	RedisServerCommandSync          = "SYNC"
	RedisServerCommandTime          = "TIME"
	RedisServerCommandAofDebug      = "AOFDEBUG"
	RedisServerCommandAofFlush      = "AOFFLUSH"
	RedisServerCommandCommand       = "COMMAND"
	RedisServerCommandExit          = "EXIT"

	RedisDebugCommandRuntimeStat = "RUNTIME_STAT"
)

func (srv *Server) doSave(cli *client.Client, encoder *rdb.Encoder) {
	loggers.Info("redis rdb save start")
	encoder.EncodeHeader()
	for dbNo, db := range srv.Databases {
		if db.DBSize() == 0 {
			continue
		}
		encoder.EncodeDatabase(dbNo)
		for _, key := range db.GetAllKeys() {
			if redisObj := db.SearchKeyInDB(key); redisObj != nil {
				switch redisObj.GetObjectType() {
				case encodings.RedisTypeString:
					if ts, ok := redisObj.(database.TString); !ok {
						cli.ResponseReError(re.ErrImpossible)
					} else {
						encoder.EncodeType(rdb.TypeString)
						encoder.EncodeRawString(key)
						encoder.EncodeRawString(ts.GetValue().(string))
					}
				case encodings.RedisTypeList:
					if tl, ok := redisObj.(database.TList); !ok {
						cli.ResponseReError(re.ErrImpossible)
					} else {
						encoder.EncodeType(rdb.TypeList)
						encoder.EncodeRawString(key)
						encoder.EncodeLength(uint32(tl.LLen()))
						for _, item := range tl.GetAllMembers() {
							encoder.EncodeRawString(item)
						}
					}
				case encodings.RedisTypeHash:
					if th, ok := redisObj.(database.THash); !ok {
						cli.ResponseReError(re.ErrImpossible)
					} else {
						encoder.EncodeType(rdb.TypeHash)
						encoder.EncodeRawString(key)
						encoder.EncodeLength(uint32(th.HLen()))
						fieldValues := th.HGetAll()
						for i := 0; i < len(fieldValues); i += 2 {
							encoder.EncodeRawString(fieldValues[i])
							encoder.EncodeRawString(fieldValues[i+1])
						}
					}
				case encodings.RedisTypeSet:
					// TODO
				case encodings.RedisTypeZSet:
					// TODO
				}
			}
		}
	}
	encoder.EncodeFooter()
	loggers.Info("redis rdb save finished")
	cli.ResponseOK()
}

// 清空Client当前所处的数据库
// TODO @lmj 这里是不是应该有锁。
func (srv *Server) FlushDB(cli *client.Client) {
	cli.SelectedDatabase().FlushDB()
	cli.ResponseOK()
}

// 清空所有数据库
func (srv *Server) FlushAll(cli *client.Client) {
	for _, db := range srv.Databases {
		db.FlushDB()
	}
	cli.ResponseOK()
}

// rdb save
func (srv *Server) Save(cli *client.Client) {
	if srv.Status.Load() == RedisServerStatusRdbSaveInProcess || srv.Status.Load() == RedisServerStatusRdbBgSaveInProcess {
		cli.ResponseReError(re.ErrRedisRdbSaveInProcess)
		return
	}
	srv.Status.Store(RedisServerStatusRdbSaveInProcess)
	defer srv.Status.Store(RedisServerStatusNormal)

	encoder, err := rdb.NewEncoder(srv.Config.RdbFilename)
	if err != nil {
		cli.ResponseReError(re.ErrUnknown)
		return
	}
	srv.doSave(cli, encoder)
}

// rdb bg save
func (srv *Server) BgSave(cli *client.Client) {
	if srv.Status.Load() == RedisServerStatusRdbSaveInProcess || srv.Status.Load() == RedisServerStatusRdbBgSaveInProcess {
		cli.ResponseReError(re.ErrRedisRdbSaveInProcess)
		return
	}
	srv.Status.Store(RedisServerStatusRdbBgSaveInProcess)
	defer srv.Status.Store(RedisServerStatusNormal)

	encoder, err := rdb.NewEncoder(srv.Config.RdbFilename)
	if err != nil {
		cli.ResponseReError(re.ErrUnknown)
		return
	}
	go srv.doSave(cli, encoder)
}

// TODO @lmj
func (srv *Server) Command(cli *client.Client) {
	cli.ResponseOK()
	loggers.Info("receive 'command' command from client, do nothing")
}

func (srv *Server) Exit(cli *client.Client) {
	cli.ResponseOK()
}

func (srv *Server) AofDebug(cli *client.Client) {
	cli.ResponseOK()
	loggers.Debug("current aof buf:%s", string(srv.aofBuf))
}

func (srv *Server) AofFlush(cli *client.Client) {
	loggers.Debug("current aof buf:%s", string(srv.aofBuf))
	srv.flushAppendOnlyFile()
	loggers.Debug("current aof buf:%s", string(srv.aofBuf))
	cli.ResponseOK()
}

func (srv *Server) RuntimeStat(cli *client.Client) {
	cli.ResponseOK()

	loggers.Info("list all clients' info ")
	for _, client := range srv.clients {
		loggers.Info("clientID:%d client:%+v", client.ID(), client)
	}
}
