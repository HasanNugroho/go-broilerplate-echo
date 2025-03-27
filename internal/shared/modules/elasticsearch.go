package modules

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type ElasticSearchConfig struct {
	Enabled  bool     `mapstructure:"ACTIVATE_ELASTICSEARCH"`
	Host     []string `mapstructure:"ELASTICSEARCH_HOST"`
	Port     int      `mapstructure:"ELASTICSEARCH_PORT"`
	Password string   `mapstructure:"ELASTICSEARCH_PASSWORD"`
	CertPath string   `mapstructure:"ELASTICSEARCH_CERT_PATH"`
	Client   *elasticsearch.Client
}

func (config *ElasticSearchConfig) SearchInit() error {
	if !config.Enabled {
		log.Warn().Msg("⚠️ Elasticsearch is disabled. Skipping initialization.")
		return nil
	}

	certPath := os.Getenv(config.CertPath)
	if certPath == "" {
		certPath = "./certs/http_ca.crt"
	}

	cert, err := os.ReadFile(certPath)
	if err != nil {
		log.Error().Msg("Error reading CA certificate: " + err.Error())
	}

	// Construct addresses with port
	var addresses []string
	for _, host := range config.Host {
		addresses = append(addresses, fmt.Sprintf("%s:%d", host, config.Port))
	}

	log.Info().Msgf("🔍 Connecting to Elasticsearch: %s", strings.Join(addresses, ", "))

	// Create Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  "elastic",
		Password:  config.Password,
		CACert:    cert,
	}

	// Initialize Elasticsearch client
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Error().Msgf("❌ Elasticsearch connection failed: %v", err)
		return fmt.Errorf("error creating the client: %w", err)
	}

	// Success!
	config.Client = es
	log.Info().Msg("✅ Elasticsearch connected successfully!")
	return nil
}
