package conf

import (
	"time"

	"github.com/SwanSpouse/redis_go/loggers"
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

	RedismaxQueryBufLen = 1024 * 1024 * 1024 /* 1GB max query buffer. */
)

// redis server configuration
type ServerConfig struct {
	/* General */
	SentinelMode int           `flag:"sentinel-mode" cfg:"sentinel-mode"` /* True if this instance is a Sentinel. */
	LogLevel     int64         `flag:"log-level" cfg:"log-level"`         /* log levels*/
	Timeout      time.Duration `flag:"timeout" cfg:"timeout"`             // Timeout represents the per-request socket read/write timeout. Default 0(disable)

	/* Networking */
	Port           int    `flag:"port" cfg:"port"`                         /* TCP listening Port */
	BindAddr       string `flag:"addr" cfg:"addr"`                         /* Bind address or NULL */
	ReaderPoolSize int    `flag:"reader-pool-size" cfg:"reader-pool-size"` /* ReaderPool 的默认大小 */
	WriterPoolSize int    `flag:"writer-pool-size" cfg:"writer-pool-size"` /* WriterPool 的默认大小*/

	/* Configuration */
	Verbosity   int   `flag:"verbosity" cfg:"verbosity"`         /* Log level in redis.conf */
	MaxIdleTime int64 `flag:"max-idle-time" cfg:"max-idle-time"` /* Client timeout in seconds */
	DBNum       int   `flag:"db-num" cfg:"db-num"`               /* Total number of configured DBs */

	ClientMaxQueryBufLen int64 `flag:"client-max-query-buf-len" cfg:"client-max-query-buf-len"`

	/* Aof persistence */
	AofState    int    `flag:"aof-state" cfg:"aof-state"`
	AofFSync    string `flag:"aof-fsync" cfg:"aof-fsync"`
	AofFilename string `flag:"aof-filename" cfg:"aof-filename"`
	//AofNoFSyncOnRewrite    int       `flag:"aof-no-fsync-on-rewrite" cfg:"aof-no-fsync-on-rewrite"`
	//AofRewritePerc         int       `flag:"aof-rewrite-perc" cfg:"aof-rewrite-perc"`
	//AofRewriteMinSize      int64     `flag:"aof-rewrite-min-size" cfg:"aof-rewrite-min-size"`
	//AofRewriteBaseSize     int64     `flag:"aof-rewrite-base-size" cfg:"aof-rewrite-base-size"`
	//AofRewriteScheduled    int       `flag:"aof-rewrite-scheduled" cfg:"aof-rewrite-scheduled"`
	//AofLastFSync           time.Time `flag:"aof-last-fsync" cfg:"aof-last-fsync"`
	//AofRewriteTimeLast     int       `flag:"aof-rewrite-time-last" cfg:"aof-rewrite-time-last"`
	//AofRewriteTimeStart    int       `flag:"aof-rewrite-time-start" cfg:"aof-rewrite-time-start"`
	//AofLastBgRewriteStatus int       `flag:"aof-last-bg-rewrite-status" cfg:"aof-last-bg-rewrite-status"`
	//AofDelayedFSync        int       `flag:"aof-delayed-fsync" cfg:"aof-delayed-fsync"`
	//AofSelectedDB          int       `flag:"aof-selected-db" cfg:"aof-selected-db"`
	//AofFlushPostponedStart int       `flag:"aof-flush-postponed-start" cfg:"aof-flush-postponed-start"`

	/* RDB persistence */
	Dirty             int64  `flag:"dirty" cfg:"dirty"`
	DirtyBeforeBgSave int64  `flag:"dirty-before-bg-save" cfg:"dirty-before-bg-save"`
	RdbFilename       string `flag:"rdb-filename" cfg:"rdb-filename"`
	//RdbCompression    int    `flag:"rdb-compression" cfg:"rdb-compression"`
	//RdbChecksum       int    `flag:"rdn-checksum" cfg:"rdn-checksum"`
	//LastSave              time.Time `flag:"last-save" cfg:"last-save"`
	//RdbSaveTimeLast       time.Time `flag:"rdb-save-time-last" cfg:"rdb-save-time-last"`
	//RdbSaveTimeStart      time.Time `flag:"rdb-save-time-start" cfg:"rdb-save-time-start"`
	//LastBgSaveStatus      int       `flag:"rdb-bg-save-status" cfg:"rdb-bg-save-status"`
	//StopWritesOnBgSaveErr int       `flag:"stop-writes-on-bg-save-err" cfg:"stop-writes-on-bg-save-err"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddr:       RedisServerAddr,
		Port:           RedisServerPort,
		DBNum:          RedisDefaultDBNum,
		LogLevel:       RedisLogLevel,
		Timeout:        5 * time.Second,
		ReaderPoolSize: RedisIOReaderPoolThreadNum,
		WriterPoolSize: RedisIOWriterPoolThreadNum,
		RdbFilename:    RedisRDBDefaultFilePath,
		AofState:       RedisAofOff,
		AofFSync:       RedisAofFSyncAlways,
		AofFilename:    RedisAofDefaultFilePath,
	}
}
