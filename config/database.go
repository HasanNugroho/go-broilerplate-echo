package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	RDBMS RDBMSConfig
	Redis RedisConfig
}

// RDBMS - relational database variables
type RDBMSConfig struct {
	Enabled bool
	Env     DBEnv
	Access  DBAccess
	Ssl     DBSsl
	Conn    DBConn
	Client  *gorm.DB
}

// REDIS - redis database variables
type RedisConfig struct {
	Enabled bool
	Env     RedisEnv
	Conn    RedisConn
	Client  *redis.Client
}

// DBEnv - environment variables for RDBMS
type DBEnv struct {
	Driver      string
	Host        string
	Port        int
	TimeZone    string
	Synchronize bool
	LogLevel    int
}

// DBAccess - database access credentials
type DBAccess struct {
	DbName string
	User   string
	Pass   string
}

// DBSsl - SSL configuration for RDBMS
type DBSsl struct {
	Mode       string
	MinTLS     string
	RootCA     string
	ServerCert string
	ClientCert string
	ClientKey  string
}

// DBConn - RDBMS connection settings
type DBConn struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisEnv - environment variables for Redis
type RedisEnv struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// RedisConn - Redis connection settings
type RedisConn struct {
	PoolSize int
	ConnTTL  int
}

// loadDatabaseConfig loads all database-related configuration
func loadDatabaseConfig() (dbConfig DatabaseConfig) {

	if activateRDBMS, err := strconv.ParseBool(os.Getenv("ACTIVATE_RDBMS")); err == nil && activateRDBMS {
		dbRDBMS := loadRDBMSConfig()
		dbConfig.RDBMS = dbRDBMS.RDBMS
		dbConfig.RDBMS.Enabled = true
	}

	if activateRedis, err := strconv.ParseBool(os.Getenv("ACTIVATE_REDIS")); err == nil && activateRedis {
		dbRedis := loadRedisConfig()
		dbConfig.Redis = dbRedis.Redis
		dbConfig.Redis.Enabled = true
	}

	return
}

func loadRDBMSConfig() (dbConfig DatabaseConfig) {
	// Load Env
	dbConfig.RDBMS.Env = DBEnv{
		Driver:      strings.ToLower(strings.TrimSpace(os.Getenv("DBDRIVER"))),
		Host:        strings.TrimSpace(os.Getenv("DBHOST")),
		Port:        utils.ToInt("DBPORT", 5432),
		TimeZone:    strings.TrimSpace(os.Getenv("DBTIMEZONE")),
		Synchronize: utils.ToBool(os.Getenv("DBSYNCHRONIZE"), false),
		LogLevel:    normalizeLogLevel(utils.ToInt("DBLOGLEVEL", 1)),
	}

	// Load Access
	dbConfig.RDBMS.Access = DBAccess{
		DbName: strings.TrimSpace(os.Getenv("DBNAME")),
		User:   strings.TrimSpace(os.Getenv("DBUSER")),
		Pass:   strings.TrimSpace(os.Getenv("DBPASS")),
	}

	// Load SSL
	dbConfig.RDBMS.Ssl = DBSsl{
		Mode:       utils.ToString(os.Getenv("DBSSLMODE"), "disable"),
		MinTLS:     strings.TrimSpace(os.Getenv("DBSSL_TLS_MIN")),
		RootCA:     strings.TrimSpace(os.Getenv("DBSSL_ROOT_CA")),
		ServerCert: strings.TrimSpace(os.Getenv("DBSSL_SERVER_CERT")),
		ClientCert: strings.TrimSpace(os.Getenv("DBSSL_CLIENT_CERT")),
		ClientKey:  strings.TrimSpace(os.Getenv("DBSSL_CLIENT_KEY")),
	}

	// Load Connection
	dbConfig.RDBMS.Conn = DBConn{
		MaxIdleConns:    utils.ToInt("DBMAXIDLECONNS", 10),
		MaxOpenConns:    utils.ToInt("DBMAXOPENCONNS", 100),
		ConnMaxLifetime: utils.ToDuration("DBCONNMAXLIFETIME", 30*time.Minute),
		ConnMaxIdleTime: utils.ToDuration("DBCONNMAXIDLETIME", 10*time.Minute),
	}

	return
}

// LoadRedisConfig loads Redis configuration
func loadRedisConfig() (dbConfig DatabaseConfig) {
	// Load Env
	dbConfig.Redis.Env = RedisEnv{
		Host:     utils.ToString(os.Getenv("REDISHOST"), "localhost"),
		Port:     utils.ToInt(os.Getenv("REDISPORT"), 6379),
		Password: utils.ToString(os.Getenv("REDISPASSWORD"), ""),
	}

	dbConfig.Redis.Conn = RedisConn{
		PoolSize: utils.ToInt(os.Getenv("POOLSIZE"), 10),
		ConnTTL:  utils.ToInt(os.Getenv("CONNTTL"), 60),
	}

	return
}

func normalizeLogLevel(logLevel int) int {
	if logLevel < 0 || logLevel > 4 {
		return 1
	}
	return logLevel
}
