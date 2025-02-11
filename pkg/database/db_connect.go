package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/HasanNugroho/starter-golang/pkg/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// dbClient variable to access gorm
var dbClient *gorm.DB

// dbClient variable to access gorm
var sqlDB *sql.DB

// redisClient to access redis
var redisClient *redis.Client

var err error

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	configureDB := config.GetConfig().Database.RDBMS

	driver := configureDB.Env.Driver
	username := configureDB.Access.User
	password := configureDB.Access.Pass
	database := configureDB.Access.DbName
	host := configureDB.Env.Host
	port := configureDB.Env.Port
	sslmode := configureDB.Ssl.Sslmode
	timeZone := configureDB.Env.TimeZone
	maxIdleConns := configureDB.Conn.MaxIdleConns
	maxOpenConns := configureDB.Conn.MaxOpenConns
	connMaxLifetime := configureDB.Conn.ConnMaxLifetime
	logLevel := configureDB.Log.LogLevel

	if logLevel < 0 || logLevel > 4 {
		logLevel = 1
	}

	switch driver {
	case "mysql":
		address := host
		if port != "" {
			address += ":" + port
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, database)
		if sslmode == "" {
			sslmode = "disable"
		}
		if sslmode != "disable" {
			if sslmode == "require" {
				dsn += "&tls=true"
			}
			if sslmode == "verify-ca" || sslmode == "verify-full" {
				dsn += "&tls=custom"
				if err := InitTLSMySQL(); err != nil {
					return nil, fmt.Errorf("❌ failed to initialize TLS for MySQL: %w", err)
				}
			}
		}

		sqlDB, err = sql.Open(driver, dsn)
		if err != nil {
			return nil, fmt.Errorf("❌ error opening MySQL connection: %w", err)
		}
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetConnMaxLifetime(connMaxLifetime)

		dbClient, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
		})
		if err != nil {
			return nil, fmt.Errorf("❌ error initializing MySQL with GORM: %w", err)
		}

	case "postgres":
		address := fmt.Sprintf("host=%s", host)
		if port != "" {
			address += fmt.Sprintf(" port=%s", port)
		}
		dsn := fmt.Sprintf("%s user=%s dbname=%s password=%s TimeZone=%s", address, username, database, password, timeZone)
		if sslmode == "" {
			sslmode = "disable"
		}
		if sslmode != "disable" {
			if configureDB.Ssl.RootCA != "" {
				dsn += fmt.Sprintf(" sslrootcert=%s", configureDB.Ssl.RootCA)
			} else if configureDB.Ssl.ServerCert != "" {
				dsn += fmt.Sprintf(" sslrootcert=%s", configureDB.Ssl.ServerCert)
			}
			if configureDB.Ssl.ClientCert != "" {
				dsn += fmt.Sprintf(" sslcert=%s", configureDB.Ssl.ClientCert)
			}
			if configureDB.Ssl.ClientKey != "" {
				dsn += fmt.Sprintf(" sslkey=%s", configureDB.Ssl.ClientKey)
			}
		}
		dsn += fmt.Sprintf(" sslmode=%s", sslmode)

		sqlDB, err = sql.Open("pgx", dsn)
		if err != nil {
			return nil, fmt.Errorf("❌ error opening PostgreSQL connection: %w", err)
		}
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetConnMaxLifetime(connMaxLifetime)

		dbClient, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
		})
		if err != nil {
			return nil, fmt.Errorf("❌ error initializing PostgreSQL with GORM: %w", err)
		}

	default:
		return nil, fmt.Errorf("the driver %s is not implemented yet", driver)
	}

	log.Println("✅ Database connection established successfully!")
	return dbClient, nil
}

// InitRedis initializes the Redis client
func InitRedis() error {
	configureRedis := config.GetConfig().Database.REDIS
	if configureRedis.Env.Host == "" || configureRedis.Env.Port == "" {
		return fmt.Errorf("❌ Redis configuration is missing host or port")
	}

	// Inisialisasi Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", configureRedis.Env.Host, configureRedis.Env.Port),
		// Password: configureRedis.Env.Password ?? "",
		// DB:       configureRedis.Env.DB ?? 0,
		PoolSize: configureRedis.Conn.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configureRedis.Conn.ConnTTL)*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	// Simpan client Redis ke konfigurasi
	configureRedis.Client = redisClient
	log.Println("✅ Redis connected successfully!")

	return nil
}
