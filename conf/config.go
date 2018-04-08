package conf

import "redis_go/log"

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
	SentinelMode int   /* True if this instance is a Sentinel. */
	LogLevel     int64 /* log levels*/

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
	sc.ReaderPoolNum = RedisIOReaderPoolThreadNum
	sc.WriterPoolNum = RedisIOWriterPoolThreadNum
	return sc
}
