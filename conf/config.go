package conf

const (
	RedisServerAddr   = ""
	RedisServerPort   = 9736
	RedisDefaultDBNum = 16
	RedisMaxIdleTime  = 0 /* default client timeout: infinite */
)

// redis server configuration
type ServerConfig struct {
	/* General */
	SentinelMode int /* True if this instance is a Sentinel. */

	/* Networking */
	Port         int    /* TCP listening Port */
	BindAddr     string /* Bind address or NULL */

	/* Configuration */
	Verbosity    int   /* Log level in redis.conf */
	MaxIdleTime  int64 /* Client timeout in seconds */
	DBNum        int   /* Total number of configured DBs */
}

func InitServerConfig() *ServerConfig {
	sc := &ServerConfig{}

	sc.BindAddr = ""
	sc.Port     = RedisServerPort
	sc.DBNum    = RedisDefaultDBNum
	return sc
}
