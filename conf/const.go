package conf

const (
	RedisConfigLocalPath   = "redis.conf"
	RedisServerDefaultPort = 65535

	// Aof
	RedisAofOpen            = "redis.aof.open"
	RedisAofFilename        = "redis.aof.filename"
	RedisAofInterval        = "redis.aof.interval"
	RedisAofFlushSize       = "redis.aof.flushsize"
	RedisAofReWriteInterval = "redis.aof.rewrite.interval"
)
