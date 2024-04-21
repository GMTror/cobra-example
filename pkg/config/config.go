package config

import (
	"time"

	"github.com/rs/xlog"
)

type Config struct {
	LogLevel    xlog.Level
	HttpPort    uint
	HttpTimeout time.Duration
	MongoURI    string
	RedisAuth   string
	RedisHost   string
	RedisPort   uint
}
