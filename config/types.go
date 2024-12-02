package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"github.com/spf13/viper"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	viper       *viper.Viper  `json:"-" yaml:"-"`
	Root        Root          `mapstructure:"root"`
	Port        uint16        `mapstructure:"port"`
	Headers     http.Header   `mapstructure:"headers"`
	Compression Compression   `mapstructure:"compression"`
	TLS         TLSConfig     `mapstructure:"tls"`
	Metrics     MetricsConfig `mapstructure:"metrics"`
}

type TLSConfig struct {
	Certificate TLSCertificate `mapstructure:"crt"`
	Key         TLSPrivateKey  `mapstructure:"key"`
	CA          TLSCertificate `mapstructure:"ca"`
	MinVersion  TLSVersion     `mapstructure:"minVersion"`
	MaxVersion  TLSVersion     `mapstructure:"maxVersion"`
	Curves      TLSCurves      `mapstructure:"curves"`
	Ciphers     TLSCiphers     `mapstructure:"ciphers"`
	ACME        ACMEConfig     `mapstructure:"acme"`
}

type ACMEConfig struct {
	Enabled       bool              `mapstructure:"enabled"`
	Email         ACMEEmail         `mapstructure:"email"`
	URL           ACMEURL           `mapstructure:"url"`
	Domains       ACMEDomains       `mapstructure:"domains"`
	ChallengePath ACMEChallengePath `mapstructure:"challengePath"`
}

type MetricsConfig struct {
	Enabled     bool
	MetricsPort uint16
}

type Root string

// TLSCertificate represents a PEM-encoded X.509 certificate.
type TLSCertificate struct {
	Block *pem.Block
	Cert  *x509.Certificate // Parsed certificate
}

// TLSPrivateKey represents a PEM-encoded private key.
type TLSPrivateKey struct {
	Block *pem.Block
	Key   interface{} // Can hold any private key type (RSA, ECDSA, ED25519, etc.)
}

type TLSVersion uint16

type TLSCurves []tls.CurveID

type TLSCiphers []uint16

type Compression int8

// ACMEEmail represents an ACME email address for registration.
type ACMEEmail string

// ACMEDomains represents the list of domains for ACME certificate issuance.
type ACMEDomains []string

// ACMEURL represents the ACME server URL.
type ACMEURL string

// ACMEChallengePath represents the path for ACME challenges.
type ACMEChallengePath string
