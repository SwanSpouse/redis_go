package conf

import "time"

// redis server configuration
type ServerConfig struct {
	Timeout time.Duration
}
