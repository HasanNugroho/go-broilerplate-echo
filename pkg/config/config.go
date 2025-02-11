package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Version  string
	Database DatabaseConfig
	// EmailConf  EmailConfig
	// Logger     LoggerConfig
	Server ServerConfig
	// Security   SecurityConfig
}

var configAll *Configuration

// Load environment variables
func InitEnv() error {
	return godotenv.Load()
}

func InitConfig() (err error) {
	err = InitEnv()
	if err != nil {
		return
	}

	var config Configuration

	// set configuration
	config.Version = strings.TrimSpace(os.Getenv("VERSION"))

	config.Database, err = database()
	if err != nil {
		return
	}

	// config.Logger = logger()

	config.Server = server()

	configAll = &config

	return
}

// GetConfig - return all the config variables
func GetConfig() *Configuration {
	return configAll
}

// database - all DB variables
func database() (databaseConfig DatabaseConfig, err error) {
	// RDBMS
	activateRDBMS, err := strconv.ParseBool(os.Getenv("ACTIVATE_RDBMS"))
	if err != nil {
		return
	}

	if activateRDBMS {
		dbRDBMS, errThis := databaseRDBMS()
		if errThis != nil {
			err = errThis
			return
		}
		databaseConfig.RDBMS = dbRDBMS.RDBMS
	}
	databaseConfig.RDBMS.Activate = activateRDBMS

	// REDIS
	activateRedis, err := strconv.ParseBool(os.Getenv("ACTIVATE_REDIS"))
	if err != nil {
		return
	}
	if activateRedis {
		dbRedis, errThis := databaseRedis()
		if errThis != nil {
			err = errThis
			return
		}
		databaseConfig.REDIS = dbRedis.REDIS
	}
	databaseConfig.REDIS.Activate = activateRedis

	return
}

func databaseRDBMS() (databaseConfig DatabaseConfig, err error) {
	// Env
	databaseConfig.RDBMS.Env.Driver = strings.ToLower(strings.TrimSpace(os.Getenv("DBDRIVER")))
	databaseConfig.RDBMS.Env.Host = strings.TrimSpace(os.Getenv("DBHOST"))
	databaseConfig.RDBMS.Env.Port = strings.TrimSpace(os.Getenv("DBPORT"))
	databaseConfig.RDBMS.Env.TimeZone = strings.TrimSpace(os.Getenv("DBTIMEZONE"))
	// Access
	databaseConfig.RDBMS.Access.DbName = strings.TrimSpace(os.Getenv("DBNAME"))
	databaseConfig.RDBMS.Access.User = strings.TrimSpace(os.Getenv("DBUSER"))
	databaseConfig.RDBMS.Access.Pass = strings.TrimSpace(os.Getenv("DBPASS"))
	// SSL
	databaseConfig.RDBMS.Ssl.Sslmode = strings.TrimSpace(os.Getenv("DBSSLMODE"))
	databaseConfig.RDBMS.Ssl.MinTLS = strings.TrimSpace(os.Getenv("DBSSL_TLS_MIN"))
	databaseConfig.RDBMS.Ssl.RootCA = strings.TrimSpace(os.Getenv("DBSSL_ROOT_CA"))
	databaseConfig.RDBMS.Ssl.ServerCert = strings.TrimSpace(os.Getenv("DBSSL_SERVER_CERT"))
	databaseConfig.RDBMS.Ssl.ClientCert = strings.TrimSpace(os.Getenv("DBSSL_CLIENT_CERT"))
	databaseConfig.RDBMS.Ssl.ClientKey = strings.TrimSpace(os.Getenv("DBSSL_CLIENT_KEY"))
	// Conn
	dbMaxIdleConns := strings.TrimSpace(os.Getenv("DBMAXIDLECONNS"))
	dbMaxOpenConns := strings.TrimSpace(os.Getenv("DBMAXOPENCONNS"))
	dbConnMaxLifetime := strings.TrimSpace(os.Getenv("DBCONNMAXLIFETIME"))
	databaseConfig.RDBMS.Conn.MaxIdleConns, err = strconv.Atoi(dbMaxIdleConns)
	if err != nil {
		return
	}
	databaseConfig.RDBMS.Conn.MaxOpenConns, err = strconv.Atoi(dbMaxOpenConns)
	if err != nil {
		return
	}
	databaseConfig.RDBMS.Conn.ConnMaxLifetime, err = time.ParseDuration(dbConnMaxLifetime)
	if err != nil {
		return
	}

	// Logger
	dbLogLevel := strings.TrimSpace(os.Getenv("DBLOGLEVEL"))
	databaseConfig.RDBMS.Log.LogLevel, err = strconv.Atoi(dbLogLevel)
	if err != nil {
		return
	}

	return
}

// databaseRedis - all REDIS DB variables
func databaseRedis() (databaseConfig DatabaseConfig, err error) {
	// REDIS
	poolSize, errThis := strconv.Atoi(strings.TrimSpace(os.Getenv("POOLSIZE")))
	if errThis != nil {
		err = errThis
		return
	}
	connTTL, errThis := strconv.Atoi(strings.TrimSpace(os.Getenv("CONNTTL")))
	if errThis != nil {
		err = errThis
		return
	}

	databaseConfig.REDIS.Env.Host = strings.TrimSpace(os.Getenv("REDISHOST"))
	databaseConfig.REDIS.Env.Port = strings.TrimSpace(os.Getenv("REDISPORT"))
	databaseConfig.REDIS.Conn.PoolSize = poolSize
	databaseConfig.REDIS.Conn.ConnTTL = connTTL

	return
}

// // logger - config for sentry.io
// func logger() (loggerConfig LoggerConfig) {
// 	loggerConfig.Activate = strings.ToLower(strings.TrimSpace(os.Getenv("ACTIVATE_SENTRY")))
// 	if loggerConfig.Activate == Activated {
// 		loggerConfig.SentryDsn = strings.TrimSpace(os.Getenv("SentryDSN"))
// 		loggerConfig.PerformanceTracing = strings.ToLower(strings.TrimSpace(os.Getenv("SENTRY_ENABLE_TRACING")))
// 		loggerConfig.TracesSampleRate = strings.TrimSpace(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"))
// 	}

// 	return
// }

// server - port and env
func server() (serverConfig ServerConfig) {
	serverConfig.ServerHost = strings.TrimSpace(os.Getenv("APP_HOST"))
	serverConfig.ServerPort = strings.TrimSpace(os.Getenv("APP_PORT"))
	serverConfig.ServerEnv = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))

	return
}
