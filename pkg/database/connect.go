package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/HasanNugroho/starter-golang/internal/core/domain"
	"github.com/HasanNugroho/starter-golang/pkg/config"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	cfg := config.GetConfig().Database.RDBMS

	logLevel := normalizeLogLevel(cfg.Env.LogLevel)

	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open(cfg.Env.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ error opening %s connection: %w", cfg.Env.Driver, err)
	}

	sqlDB.SetMaxIdleConns(cfg.Conn.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Conn.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Conn.ConnMaxLifetime)

	gormDB, err := gorm.Open(getGormDialector(cfg.Env.Driver, sqlDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ error initializing %s with GORM: %w", cfg.Env.Driver, err)
	}

	log.Println("✅ Database connection established successfully!")

	if cfg.Env.Synchronize {
		log.Println("Running AutoMigrate...")
		if err := gormDB.AutoMigrate(&domain.User{}); err != nil {
			return nil, fmt.Errorf("❌ migration failed: %w", err)
		}
		log.Println("✅ Successfully Migrated")
	}

	return gormDB, nil
}

// InitRedis initializes the Redis client
func InitRedis() (*redis.Client, error) {
	cfg := config.GetConfig().Database.REDIS

	if cfg.Env.Host == "" || cfg.Env.Port == "" {
		return nil, fmt.Errorf("❌ Redis configuration is missing host or port")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Env.Host, cfg.Env.Port),
		PoolSize: cfg.Conn.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Conn.ConnTTL)*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	log.Println("✅ Redis connected successfully!")
	return redisClient, nil
}

// normalizeLogLevel ensures log level is within range
func normalizeLogLevel(logLevel int) int {
	if logLevel < 0 || logLevel > 4 {
		return 1
	}
	return logLevel
}

// getDSN builds the DSN string for MySQL or PostgreSQL
func getDSN(cfg config.RDBMS) (string, error) {
	switch cfg.Env.Driver {
	case "mysql":
		return buildMySQLDSN(cfg)
	case "postgres":
		return buildPostgresDSN(cfg), nil
	default:
		return "", fmt.Errorf("❌ the driver %s is not implemented yet", cfg.Env.Driver)
	}
}

// buildMySQLDSN constructs MySQL DSN
func buildMySQLDSN(cfg config.RDBMS) (string, error) {
	address := cfg.Env.Host
	if cfg.Env.Port != "" {
		address += ":" + cfg.Env.Port
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Access.User, cfg.Access.Pass, address, cfg.Access.DbName)

	if cfg.Ssl.Sslmode == "" {
		cfg.Ssl.Sslmode = "disable"
	}

	if cfg.Ssl.Sslmode != "disable" {
		switch cfg.Ssl.Sslmode {
		case "require":
			dsn += "&tls=true"
		case "verify-ca", "verify-full":
			dsn += "&tls=custom"
			if err := InitTLSMySQL(); err != nil {
				return "", fmt.Errorf("❌ failed to initialize TLS for MySQL: %w", err)
			}
		}
	}

	return dsn, nil
}

// buildPostgresDSN constructs PostgreSQL DSN
func buildPostgresDSN(cfg config.RDBMS) string {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=%s",
		cfg.Env.Host, cfg.Access.User, cfg.Access.DbName, cfg.Access.Pass, cfg.Env.TimeZone, cfg.Ssl.Sslmode)

	if cfg.Ssl.RootCA != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", cfg.Ssl.RootCA)
	} else if cfg.Ssl.ServerCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", cfg.Ssl.ServerCert)
	}
	if cfg.Ssl.ClientCert != "" {
		dsn += fmt.Sprintf(" sslcert=%s", cfg.Ssl.ClientCert)
	}
	if cfg.Ssl.ClientKey != "" {
		dsn += fmt.Sprintf(" sslkey=%s", cfg.Ssl.ClientKey)
	}

	return dsn
}

// getGormDialector returns the correct GORM dialector
func getGormDialector(driver string, sqlDB *sql.DB) gorm.Dialector {
	switch driver {
	case "mysql":
		return mysql.New(mysql.Config{Conn: sqlDB})
	case "postgres":
		return postgres.New(postgres.Config{Conn: sqlDB})
	default:
		panic(fmt.Sprintf("❌ Unsupported driver: %s", driver))
	}
}

// ShutdownDB closes the database connection using GORM
func ShutdownDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("❌ Error retrieving database connection: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("❌ Error closing database connection: %v", err)
		} else {
			log.Println("✅ Database connection closed successfully!")
		}
	}
}

// ShutdownRedis closes the Redis connection
func ShutdownRedis(redisClient *redis.Client) {
	if redisClient != nil {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redisClient.Close(); err != nil {
			log.Printf("❌ Error closing Redis connection: %v", err)
		} else {
			log.Println("✅ Redis connection closed successfully!")
		}
	}
}
