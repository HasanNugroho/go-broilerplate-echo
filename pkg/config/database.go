package config

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type DatabaseConfig struct {
	// relational database
	RDBMS RDBMS

	// redis database
	REDIS REDIS
}

// RDBMS - relational database variables
type RDBMS struct {
	Activate bool
	Env      struct {
		Driver      string
		Host        string
		Port        string
		TimeZone    string
		Synchronize bool
		LogLevel    int
	}
	Access struct {
		DbName string
		User   string
		Pass   string
	}
	Ssl struct {
		Sslmode    string
		MinTLS     string
		RootCA     string
		ServerCert string
		ClientCert string
		ClientKey  string
	}
	Conn struct {
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
	}
}

// REDIS - redis database variables
type REDIS struct {
	Activate bool
	Env      struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
	Conn struct {
		PoolSize int
		ConnTTL  int
	}
	Client *redis.Client
}
