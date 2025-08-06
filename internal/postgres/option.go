package postgres

import "time"

type ConnectionOptions struct {
	URL            string
	MaxConnections int
	MaxIdleTime    time.Duration
	ConnectTimeout time.Duration
}
