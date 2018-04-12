package conf

import (
	"redis_go/log"
	"time"
)

const (
	RedisServerAddr   = ""
	RedisServerPort   = 9736
	RedisDefaultDBNum = 16
	RedisMaxIdleTime  = 0 /* default client timeout: infinite */

	RedisIOReaderPoolThreadNum = 5
	RedisIOWriterPoolThreadNum = 5

	RedisLogLevel = log.DEBUG
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
	return sc
}
