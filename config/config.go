package config

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Build initializes the Config by loading from environment variables and optionally from a YAML file.
func Build() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get cwd: %w", err)
	}
	cfg := &Config{
		viper:       viper.NewWithOptions(viper.EnvKeyReplacer(&envReplacer{})),
		Root:        Root(cwd),
		Port:        8080,
		Headers:     http.Header{},
		Compression: Compression(gzip.BestSpeed),
		Metrics: MetricsConfig{
			Enabled:     false,
			MetricsPort: 9090,
		},
		TLS: TLSConfig{
			MinVersion: TLSVersion(tls.VersionTLS12),
			Curves: TLSCurves{
				tls.CurveP256,
				tls.CurveP384,
				tls.CurveP521,
			},
			Ciphers: TLSCiphers{
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			},
			ACME: ACMEConfig{
				Enabled: false,
			},
		},
	}

	// Initialize Viper
	cfg.viper.SetConfigType("yaml")   // Set the config type to YAML
	cfg.viper.AddConfigPath(".")      // Add the current directory as the config path
	cfg.viper.SetConfigName("config") // The name of the config file (config.yaml)
	// Get the environment variable prefix from GSS_ENV_PREFIX
	envPrefix := os.Getenv("GSS_ENV_PREFIX")
	configPathEnvName := "CONFIG_PATH"
	if envPrefix != "" {
		cfg.viper.SetEnvPrefix(envPrefix) // Set the environment variable prefix dynamically
		configPathEnvName = strings.ToUpper(fmt.Sprintf("%s_%s", envPrefix, configPathEnvName))
	}
	// Check if CONFIG_PATH is set to load config from a YAML file
	configPath := os.Getenv(configPathEnvName)
	if configPath != "" {
		cfg.viper.SetConfigFile(configPath) // Set the config file path
		// Read the config file
		if err := cfg.viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("unable to read config file: %w", err)
		}
	}

	// Automatically read environment variables
	cfg.viper.AutomaticEnv()

	// Unmarshal the configuration into the struct
	if err := cfg.viper.Unmarshal(cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			DecodeCompression, DecodeHeaders, DecodeRoot,
			DecodeTLSCertificate, DecodeTLSPrivateKey,
			DecodeTLSCurves, DecodeTLSCiphers, DecodeTLSVersion,
			DecodeACMEDomains, DecodeACMEURL, DecodeACMEChallengePath, DecodeACMEEmail,
		),
	)); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config into struct: %w", err)
	}
	// Return nil if everything went fine
	return cfg, nil
}

type envReplacer struct{}

func (r *envReplacer) Replace(s string) string {
	return strings.ToUpper(strings.NewReplacer(".", "_").Replace(s))
}
