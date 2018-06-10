package conf

import (
	"redis_go/loggers"
	"time"
)

const (
	RedisServerAddr   = ""
	RedisServerPort   = 9736
	RedisDefaultDBNum = 16
	RedisMaxIdleTime  = 0 /* default client timeout: infinite */

	RedisIOReaderPoolThreadNum = 5
	RedisIOWriterPoolThreadNum = 5

	RedisLogLevel = loggers.DEBUG

	/* AOF states */
	RedisAofOff             = 0 /* AOF is off */
	RedisAofOn              = 1 /* AOF is on */
	RedisAofWaitRewrite     = 2 /* AOF waits rewrite to start appending */
	RedisAofAppendOnlyYes   = "yes"
	RedisAofAppendonlyNo    = "no"
	RedisAofFSyncAlways     = "always"
	RedisAofFSyncEverySec   = "everysec"
	RedisAofFSyncNo         = "no"
	RedisAofDefaultFilePath = "appendonly.aof"

	/* RDB persistence */
	RedisRDBDefaultFilePath = "dump.rdb"
)

// redis server configuration
type ServerConfig struct {
	/* General */
	SentinelMode int           /* True if this instance is a Sentinel. */
	LogLevel     int64         /* log levels*/
	Timeout      time.Duration // Timeout represents the per-request socket read/write timeout. Default 0(disable)

	/* Networking */
	Port          int    /* TCP listening Port */
	BindAddr      string /* Bind address or NULL */
	ReaderPoolNum int    /* ReaderPool 的默认大小 */
	WriterPoolNum int    /* WriterPool 的默认大小*/

	/* Configuration */
	Verbosity   int   /* Log level in redis.conf */
	MaxIdleTime int64 /* Client timeout in seconds */
	DBNum       int   /* Total number of configured DBs */

	/* Aof persistence */
	AofState    int
	AofFSync    string
	AofFilename string

	/* RDB persistence */
	Dirty                 int64
	DirtyBeforeBgSave     int64
	RdbFilename           string
	RdbCompression        int
	RdbChecksum           int
	LastSave              time.Time
	RdbSaveTimeLast       time.Time
	RdbSaveTimeStart      time.Time
	LastBgSaveStatus      int
	StopWritesOnBgSaveErr int
}

func InitServerConfig() *ServerConfig {
	sc := &ServerConfig{}

	sc.BindAddr = ""
	sc.Port = RedisServerPort
	sc.DBNum = RedisDefaultDBNum
	sc.LogLevel = RedisLogLevel
	sc.Timeout = 5 * time.Second
	sc.ReaderPoolNum = RedisIOReaderPoolThreadNum
	sc.WriterPoolNum = RedisIOWriterPoolThreadNum
	sc.RdbFilename = RedisRDBDefaultFilePath
	sc.AofState = RedisAofOff
	sc.AofFSync = RedisAofFSyncAlways
	sc.AofFilename = RedisAofDefaultFilePath
	return sc
}
