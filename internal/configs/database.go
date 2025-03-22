package configs

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

	"github.com/HasanNugroho/starter-golang/internal/pkg/utils"
	driver_mysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBClient *gorm.DB

// RDBMS - relational database variables
type RDBMSConfig struct {
	Enabled         bool
	Driver          string
	Host            string
	Port            int
	DB              string
	User            string
	Pass            string
	TimeZone        string
	LogLevel        int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Ssl             DBSsl
	Client          *gorm.DB
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

func loadRDBMSConfig() (dbConfig RDBMSConfig) {
	dbConfig = RDBMSConfig{
		Enabled:         utils.ToBool(os.Getenv("ACTIVATE_RDBMS"), false),
		Driver:          strings.ToLower(utils.ToString(os.Getenv("DBDRIVER"), "postgres")),
		Host:            utils.ToString(os.Getenv("DBHOST"), "localhost"),
		Port:            utils.ToInt(os.Getenv("DBPORT"), 5432),
		DB:              utils.ToString(os.Getenv("DBNAME"), ""),
		User:            utils.ToString(os.Getenv("DBUSER"), ""),
		Pass:            utils.ToString(os.Getenv("DBPASS"), ""),
		TimeZone:        utils.ToString(os.Getenv("DBTIMEZONE"), "Asia/Jakarta"),
		LogLevel:        normalizeLogLevel(utils.ToInt(os.Getenv("DBLOGLEVEL"), 1)),
		MaxIdleConns:    utils.ToInt(os.Getenv("DBMAXIDLECONNS"), 10),
		MaxOpenConns:    utils.ToInt(os.Getenv("DBMAXOPENCONNS"), 100),
		ConnMaxLifetime: utils.ToDuration(os.Getenv("DBCONNMAXLIFETIME"), 30*time.Minute),
		ConnMaxIdleTime: utils.ToDuration(os.Getenv("DBCONNMAXIDLETIME"), 10*time.Minute),
	}

	// Load SSL
	dbConfig.Ssl = DBSsl{
		Mode:       utils.ToString(os.Getenv("DBSSLMODE"), "disable"),
		MinTLS:     utils.ToString(os.Getenv("DBSSL_TLS_MIN"), ""),
		RootCA:     utils.ToString(os.Getenv("DBSSL_ROOT_CA"), ""),
		ServerCert: utils.ToString(os.Getenv("DBSSL_SERVER_CERT"), ""),
		ClientCert: utils.ToString(os.Getenv("DBSSL_CLIENT_CERT"), ""),
		ClientKey:  utils.ToString(os.Getenv("DBSSL_CLIENT_KEY"), ""),
	}

	return
}

const (
	maxRetries = 5
	retryDelay = 3 * time.Second
)

// InitDB initializes the database connection
func InitDB(cfg *RDBMSConfig) (*gorm.DB, error) {
	// Get GORM dialect based on the database driver
	dialect, err := getDialect(cfg)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to get dialect for %s: %w", cfg.Driver, err)
	}

	var sqlDB *sql.DB
	var db *gorm.DB
	for i := 1; i <= maxRetries; i++ {
		// Open database connection
		db, err = gorm.Open(dialect, &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
		})

		if err == nil {
			Logger.Info().Msgf("✅ Database %s connected successfully!", cfg.Driver)
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
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	cfg.Client = db

	return db, nil
}

func getDialect(cfg *RDBMSConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "mysql":
		mysqlDsn, err := buildMySQLDSN(cfg)
		if err != nil {
			return nil, fmt.Errorf("❌ Build dsn of %s failed", cfg.Driver)
		}
		Logger.Error().Msg(mysqlDsn)
		return mysql.Open(mysqlDsn), nil
	case "postgres":
		return postgres.Open(buildPostgresDSN(cfg)), nil
	default:
		return nil, fmt.Errorf("❌ the driver %s is not implemented yet", cfg.Driver)
	}
}

// buildMySQLDSN constructs MySQL DSN
func buildMySQLDSN(cfg *RDBMSConfig) (string, error) {
	address := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	Logger.Error().Msg(strconv.Itoa(cfg.Port))

	// Format dasar DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, address, cfg.DB,
	)

	// Konfigurasi TLS jika diaktifkan
	if cfg.Ssl.Mode != "disable" {
		switch cfg.Ssl.Mode {
		case "require":
			dsn += "&tls=true"
		case "verify-ca", "verify-full":
			// Inisialisasi TLS lebih dulu sebelum mengubah DSN
			if err := InitTLSMySQL(cfg); err != nil {
				return "", fmt.Errorf("❌ failed to initialize TLS for MySQL: %w", err)
			}
			dsn += "&tls=custom"
		}
	}

	return dsn, nil
}

// buildPostgresDSN constructs PostgreSQL DSN
func buildPostgresDSN(cfg *RDBMSConfig) string {
	var dsn strings.Builder

	// Tambahkan parameter dasar
	dsn.WriteString(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.DB, cfg.Pass, cfg.TimeZone, cfg.Ssl.Mode,
	))

	if cfg.Ssl.Mode != "disable" {
		// Tambahkan sertifikat SSL jika tersedia
		if cfg.Ssl.RootCA != "" {
			dsn.WriteString(fmt.Sprintf(" sslrootcert=%s", cfg.Ssl.RootCA))
		} else if cfg.Ssl.ServerCert != "" {
			dsn.WriteString(fmt.Sprintf(" sslrootcert=%s", cfg.Ssl.ServerCert))
		}

		if cfg.Ssl.ClientCert != "" {
			dsn.WriteString(fmt.Sprintf(" sslcert=%s", cfg.Ssl.ClientCert))
		}
		if cfg.Ssl.ClientKey != "" {
			dsn.WriteString(fmt.Sprintf(" sslkey=%s", cfg.Ssl.ClientKey))
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
func InitTLSMySQL(cfg *RDBMSConfig) error {
	rootCA := cfg.Ssl.RootCA
	serverCert := cfg.Ssl.ServerCert
	clientCert := cfg.Ssl.ClientCert
	clientKey := cfg.Ssl.ClientKey
	minTLS := cfg.Ssl.MinTLS

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
