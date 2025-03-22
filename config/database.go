package config

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	driver_mysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBClient *gorm.DB

const (
	maxRetries = 5
	retryDelay = 3 * time.Second
)

// InitDB initializes the database connection
func InitDB(config *DBConfig) (*gorm.DB, error) {
	// Get GORM dialect based on the database driver
	dialect, err := getDialect(config)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to get dialect for %s: %w", config.Driver, err)
	}

	var sqlDB *sql.DB
	var db *gorm.DB
	for i := 1; i <= maxRetries; i++ {
		// Open database connection
		db, err = gorm.Open(dialect, &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(config.LogLevel)),
		})

		if err == nil {
			Logger.Info().Msgf("✅ Database %s connected successfully!", config.Driver)
			break
		}

		Logger.Error().Msgf("❌ Database connection failed: %v", err)
		if i < maxRetries {
			Logger.Error().Msgf("🔄 Retrying in %v seconds... (%d/%d)", retryDelay.Seconds(), i, maxRetries)
			time.Sleep(retryDelay)
		} else {
			return nil, fmt.Errorf("❌ Failed to connect after %d attempts: %v", maxRetries, err)
		}
	}

	// Retrieve the underlying SQL DB connection
	sqlDB, err = db.DB()
	if err != nil {
		return nil, fmt.Errorf("❌ error retrieving SQL DB from GORM: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	config.Client = db

	return db, nil
}

func getDialect(config *DBConfig) (gorm.Dialector, error) {
	switch config.Driver {
	case "mysql":
		mysqlDsn, err := buildMySQLDSN(config)
		if err != nil {
			return nil, fmt.Errorf("❌ Build dsn of %s failed", config.Driver)
		}
		Logger.Error().Msg(mysqlDsn)
		return mysql.Open(mysqlDsn), nil
	case "postgres":
		return postgres.Open(buildPostgresDSN(config)), nil
	default:
		return nil, fmt.Errorf("❌ the driver %s is not implemented yet", config.Driver)
	}
}

// buildMySQLDSN constructs MySQL DSN
func buildMySQLDSN(config *DBConfig) (string, error) {
	address := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	Logger.Error().Msg(strconv.Itoa(config.Port))

	// Format dasar DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Pass, address, config.DB,
	)

	// Konfigurasi TLS jika diaktifkan
	if config.Ssl.Mode != "disable" {
		switch config.Ssl.Mode {
		case "require":
			dsn += "&tls=true"
		case "verify-ca", "verify-full":
			// Inisialisasi TLS lebih dulu sebelum mengubah DSN
			if err := InitTLSMySQL(config); err != nil {
				return "", fmt.Errorf("❌ failed to initialize TLS for MySQL: %w", err)
			}
			dsn += "&tls=custom"
		}
	}

	return dsn, nil
}

// buildPostgresDSN constructs PostgreSQL DSN
func buildPostgresDSN(config *DBConfig) string {
	var dsn strings.Builder

	// Tambahkan parameter dasar
	dsn.WriteString(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=%s",
		config.Host, config.User, config.DB, config.Pass, config.TimeZone, config.Ssl.Mode,
	))

	if config.Ssl.Mode != "disable" {
		// Tambahkan sertifikat SSL jika tersedia
		if config.Ssl.RootCA != "" {
			dsn.WriteString(fmt.Sprintf(" sslrootcert=%s", config.Ssl.RootCA))
		} else if config.Ssl.ServerCert != "" {
			dsn.WriteString(fmt.Sprintf(" sslrootcert=%s", config.Ssl.ServerCert))
		}

		if config.Ssl.ClientCert != "" {
			dsn.WriteString(fmt.Sprintf(" sslcert=%s", config.Ssl.ClientCert))
		}
		if config.Ssl.ClientKey != "" {
			dsn.WriteString(fmt.Sprintf(" sslkey=%s", config.Ssl.ClientKey))
		}
	}

	return dsn.String()
}

// ShutdownDB closes the database connection using GORM
func ShutdownDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			Logger.Error().Msgf("❌ Error retrieving database connection: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			Logger.Error().Msgf("❌ Error closing database connection: %v", err)
		} else {
			Logger.Info().Msg("✅ Database connection closed successfully!")
		}
	}
}

// InitTLSMySQL initializes TLS configuration for MySQL
func InitTLSMySQL(config *DBConfig) error {
	rootCA := config.Ssl.RootCA
	serverCert := config.Ssl.ServerCert
	clientCert := config.Ssl.ClientCert
	clientKey := config.Ssl.ClientKey
	minTLS := config.Ssl.MinTLS

	// Load Root CA
	rootCertPool := x509.NewCertPool()
	var pem []byte
	var err error

	if rootCA != "" {
		pem, err = os.ReadFile(rootCA)
		if err != nil {
			return fmt.Errorf("failed to read Root CA file: %w", err)
		}
	} else {
		if serverCert == "" {
			return errors.New("missing Root CA and server certificate")
		}

		pem, err = os.ReadFile(serverCert)
		if err != nil {
			return fmt.Errorf("failed to read server certificate: %w", err)
		}
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return errors.New("failed to parse PEM encoded certificate")
	}

	// Define TLS Configuration
	tlsConfig := tls.Config{
		RootCAs: rootCertPool,
	}

	// Set Minimum TLS Version
	switch minTLS {
	case "1.1":
		tlsConfig.MinVersion = tls.VersionTLS11
	case "1.2":
		tlsConfig.MinVersion = tls.VersionTLS12
	case "1.3":
		tlsConfig.MinVersion = tls.VersionTLS13
	default:
		tlsConfig.MinVersion = tls.VersionTLS12 // Default TLS 1.2
	}

	// Load Client Certificate & Key (if provided)
	if clientCert != "" && clientKey != "" {
		cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return fmt.Errorf("failed to load client certificate & key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Register TLS Config for MySQL Driver
	err = driver_mysql.RegisterTLSConfig("custom", &tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to register custom TLS config: %w", err)
	}

	Logger.Info().Msg("✅ MySQL TLS configuration successfully registered!")
	return nil
}

func normalizeLogLevel(logLevel int) int {
	if logLevel < 0 || logLevel > 4 {
		return 1
	}
	return logLevel
}
