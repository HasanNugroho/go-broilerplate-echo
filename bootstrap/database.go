package bootstrap

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/HasanNugroho/starter-golang/config"
	driver_mysql "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.RDBMSConfig) (*gorm.DB, error) {
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
	sqlDB.SetConnMaxIdleTime(cfg.Conn.ConnMaxIdleTime)

	gormDB, err := gorm.Open(getGormDialector(cfg.Env.Driver, sqlDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.Env.LogLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ error initializing %s with GORM: %w", cfg.Env.Driver, err)
	}

	// if cfg.Env.Synchronize {
	// 	log.Println("Running AutoMigrate...")
	// 	if err := gormDB.AutoMigrate(&domain.User{}); err != nil {
	// 		return nil, fmt.Errorf("❌ migration failed: %w", err)
	// 	}
	// 	log.Println("✅ Successfully Migrated")
	// }

	cfg.Client = gormDB

	return gormDB, nil
}

// getDSN builds the DSN string for MySQL or PostgreSQL
func getDSN(cfg *config.RDBMSConfig) (string, error) {
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
func buildMySQLDSN(cfg *config.RDBMSConfig) (string, error) {
	address := cfg.Env.Host
	if strconv.Itoa(cfg.Env.Port) != "" {
		address += ":" + strconv.Itoa(cfg.Env.Port)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Access.User, cfg.Access.Pass, address, cfg.Access.DbName)

	if cfg.Ssl.Mode != "disable" {
		switch cfg.Ssl.Mode {
		case "require":
			dsn += "&tls=true"
		case "verify-ca", "verify-full":
			dsn += "&tls=custom"
			if err := InitTLSMySQL(cfg); err != nil {
				return "", fmt.Errorf("❌ failed to initialize TLS for MySQL: %w", err)
			}
		}
	}

	return dsn, nil
}

// buildPostgresDSN constructs PostgreSQL DSN
func buildPostgresDSN(cfg *config.RDBMSConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=%s",
		cfg.Env.Host, cfg.Access.User, cfg.Access.DbName, cfg.Access.Pass, cfg.Env.TimeZone, cfg.Ssl.Mode)

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

// InitTLSMySQL registers a custom tls.Config
//
// Tutorial: How to configure MySQL instance and enable TLS support
//
// 1.0 generate CA's private key and certificate
//
// to omit password: `-nodes -keyout`
//
// `openssl req -x509 -sha512 -newkey rsa:4096 -days 10950 -keyout ca-key.pem -out ca.pem`
//
// 2.0 generate web server's private key and certificate signing request (CSR)
//
// to omit password: `-nodes -keyout`
//
// Common Name (e.g. server FQDN or YOUR name) must be different for CA and web server certificates
//
// `openssl req -sha512 -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem`
//
// 2.1 config file
//
// IP: server's public or local IPs of the interfaces
//
// `echo "subjectAltName=DNS:localhost,IP:127.0.0.1,IP:172.17.0.1,IP:x.x.x.x,IP:y.y.y.y" > "server-ext.cnf"`
//
// 2.2 use CA's private key to sign web server's CSR and get back the signed certificate
//
// `openssl x509 -sha512 -req -in server-req.pem -days 3650 -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf`
//
// 2.3 verify the certificate
//
// `openssl verify -CAfile ca.pem server-cert.pem`
//
// 3.0 convert PKCS#8 format key into PKCS#1 format
//
// `openssl rsa -in server-key.pem -out server-key.pem`
//
// 4.0 replace existing files located at /var/lib/mysql
//
// 5.0 set ownership and r/w permissions
//
// ```bash
// sudo chown -R mysql:mysql ca-key.pem ca.pem server-key.pem server-cert.pem
// sudo chmod -R 600 ca-key.pem server-key.pem
// sudo chmod -R 644 ca.pem server-cert.pem
// ```
//
// 6.0 restart mysql service
//
// `sudo service mysql restart`
//
// 7.0 optional:
//
// 7.1 generate client's private key and certificate signing request (CSR)
//
// `openssl req -sha512 -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem`
//
// 7.2 config file
//
// IP: server's public or local IPs of the interfaces
//
// `echo "subjectAltName=DNS:localhost,IP:127.0.0.1,IP:172.17.0.1,IP:x.x.x.x,IP:y.y.y.y" > "client-ext.cnf"`
//
// 7.3 use CA's private key to sign client's CSR and get back the signed certificate
//
// `openssl x509 -sha512 -req -in client-req.pem -days 3650 -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf`
//
// 7.4 verify the certificate
//
// `openssl verify -CAfile ca.pem client-cert.pem`
//
// 7.5 convert PKCS#8 format key into PKCS#1 format
//
// `openssl rsa -in client-key.pem -out client-key.pem`
func InitTLSMySQL(cfg *config.RDBMSConfig) (err error) {
	minTLS := cfg.Ssl.MinTLS
	rootCA := cfg.Ssl.RootCA
	serverCert := cfg.Ssl.ServerCert
	clientCert := cfg.Ssl.ClientCert
	clientKey := cfg.Ssl.ClientKey

	rootCertPool := x509.NewCertPool()
	var pem []byte

	if rootCA != "" {
		pem, err = os.ReadFile(rootCA)
		if err != nil {
			return
		}
	} else {
		if serverCert == "" {
			err = errors.New("missing server certificate")
			return
		}

		pem, err = os.ReadFile(serverCert)
		if err != nil {
			return
		}
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		err = errors.New("failed to parse PEM encoded certificates")
		return
	}

	tlsConfig := tls.Config{}
	tlsConfig.MinVersion = tls.VersionTLS12 // default: TLS 1.2

	if minTLS == "1.1" {
		tlsConfig.MinVersion = tls.VersionTLS11
	}
	if minTLS == "1.2" {
		tlsConfig.MinVersion = tls.VersionTLS12
	}
	if minTLS == "1.3" {
		tlsConfig.MinVersion = tls.VersionTLS13
	}
	tlsConfig.RootCAs = rootCertPool

	if clientCert != "" && clientKey != "" {
		clientCertificate := make([]tls.Certificate, 0, 1)
		var certs tls.Certificate

		certs, err = tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return
		}

		clientCertificate = append(clientCertificate, certs)
		tlsConfig.Certificates = clientCertificate
	}

	err = driver_mysql.RegisterTLSConfig("custom", &tlsConfig)

	return
}
