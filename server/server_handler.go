package server

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/rdb"
)

var (
	_ client.BaseHandler = (*Server)(nil)
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
)

func (srv *Server) Process(cli *client.Client) {
	switch cli.GetCommandName() {
	case RedisServerCommandBGSRewriteAof:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandBGSave:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandClient:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandConfig:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandDBSize:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandDebug:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandFlushAll:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandFlushDB:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandInfo:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandLastSave:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandMonitor:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandPSync:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandSave:
		srv.commandSave(cli)
	case RedisServerCommandShutDown:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandSlaveOf:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandSlowLog:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandSync:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisServerCommandTime:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.GetOriginCommandName())
	}
	cli.Flush()
}

// rdb save
func (srv *Server) commandSave(cli *client.Client) {
	encoder, err := rdb.NewEncoder("dump.rdb.test")
	if err != nil {
		cli.ResponseReError(re.ErrUnknown)
		return
	}
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
	cli.ResponseOK()
}
