package config_test

import (
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vlasov-y/gss/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var (
	cfg         *config.Config
	ok          bool
	err         error
	tempDir     string
	configFile  string
	certFile    string
	keyFile     string
	certContent string
	keyContent  string
)

// Generate certificate and key before tests
var _ = BeforeEach(func() {
	// Make a temp dir for test files
	tempDir, err = os.MkdirTemp("", "gss-*")
	Expect(err).ToNot(HaveOccurred())
	lstat, err := os.Lstat(tempDir)
	Expect(err).ToNot(HaveOccurred())
	Expect(lstat.IsDir()).To(BeTrue())

	configFile = filepath.Join(tempDir, "config.yaml")
	certFile = filepath.Join(tempDir, "cert.pem")
	keyFile = filepath.Join(tempDir, "key.pem")

	// Generate a self-signed certificate and key
	Expect(generateTLSCertificate(certFile, keyFile)).To(Succeed())

	// Read certificate and key and save to variables
	var c, k []byte
	c, err = os.ReadFile(certFile)
	Expect(err).ToNot(HaveOccurred())
	k, err = os.ReadFile(keyFile)
	Expect(err).ToNot(HaveOccurred())
	certContent, keyContent = string(c), string(k)

	// Initialize the config
	cfg, err = config.Build()
	Expect(err).ToNot(HaveOccurred())
})

// Clean up certificate and key after tests
var _ = AfterEach(func() {
	if tempDir != "" {
		err = os.RemoveAll(tempDir)
		Expect(err).ToNot(HaveOccurred())
	}
})

var _ = Describe("Configuration management", func() {
	It("should load configuration from a YAML file and overwrite from ENV", func() {
		// Create a sample Config object
		yamlData := fmt.Sprintf(`
port: 8888
metrics:
  enabled: true
headers:
  key: value
compression: default
tls:
  crt: %s
  key: %s
  ca: %s
  minVersion: TLS1.1
  curves:
    - P-521
  ciphers:
    - TLS_CHACHA20_POLY1305_SHA256
  acme:
    enabled: true
    email: admin@example.com
    url: https://example.com/
    domains:
      - example.com
      - www.example.com
    challengePath: /.well-known/acme-challenge/
`, fmt.Sprintf("%q", certContent), fmt.Sprintf("%q", keyContent), certFile)

		// Write the YAML content to the temporary file
		err = os.WriteFile(configFile, []byte(yamlData), 0o644)
		Expect(err).ToNot(HaveOccurred())

		// Env variables override
		env := map[string]string{
			"GSS_ENV_PREFIX":       "gss",
			"GSS_ROOT":             "/tmp/../tmp",
			"GSS_CONFIG_PATH":      configFile,
			"GSS_TLS_CIPHERS":      "TLS_RSA_WITH_AES_128_CBC_SHA256",
			"GSS_TLS_ACME_ENABLED": "true",
			"GSS_TLS_MAXVERSION":   "TLS1.2",
		}
		for name, value := range env {
			Expect(os.Setenv(name, value)).To(Succeed())
		}

		// Initialize the config from YAML file
		cfg, err = config.Build()
		Expect(err).ToNot(HaveOccurred())

		for name := range env {
			Expect(os.Unsetenv(name)).To(Succeed())
		}

		// Verify that the values in the config match the values in the YAML file
		Expect(cfg.Root).To(Equal(config.Root("/tmp")))
		Expect(cfg.Port).To(Equal(uint16(8888)))
		Expect(cfg.Metrics.MetricsPort).To(Equal(uint16(9090)))
		Expect(cfg.Metrics.Enabled).To(BeTrue())

		// Verify Headers and Compression
		Expect(cfg.Headers).ToNot(BeEmpty())
		headers := http.Header{}
		headers.Add("key", "value")
		Expect(cfg.Headers).To(Equal(headers))
		Expect(cfg.Compression).To(Equal(config.Compression(gzip.DefaultCompression)))

		// Verify TLSConfig values
		Expect(cfg.TLS.Certificate.Block).ToNot(BeNil())
		Expect(cfg.TLS.Key.Block).ToNot(BeNil())
		Expect(cfg.TLS.CA.Block).ToNot(BeNil())
		Expect(cfg.TLS.MinVersion).To(Equal(config.TLSVersion(tls.VersionTLS11)))
		Expect(cfg.TLS.MaxVersion).To(Equal(config.TLSVersion(tls.VersionTLS12)))

		// Verify TLSCurves
		Expect(cfg.TLS.Curves).To(ContainElement(tls.CurveP521))

		// Verify TLSCiphers
		Expect(cfg.TLS.Ciphers).To(ContainElement(tls.TLS_RSA_WITH_AES_128_CBC_SHA256))

		// Verify ACMEConfig values
		Expect(cfg.TLS.ACME.Enabled).To(BeTrue())
		Expect(cfg.TLS.ACME.Email).To(Equal(config.ACMEEmail("admin@example.com")))
		Expect(cfg.TLS.ACME.URL).To(Equal(config.ACMEURL("https://example.com/")))
		Expect(cfg.TLS.ACME.Domains).To(ContainElement("example.com"))
		Expect(cfg.TLS.ACME.Domains).To(ContainElement("www.example.com"))
		Expect(cfg.TLS.ACME.ChallengePath).To(Equal(config.ACMEChallengePath("/.well-known/acme-challenge/")))
	})

	It("should return an error for a broken config", func() {
		// Create a sample Config object
		yamlData := `
headers:
  key:
		sub: key
`

		// Write the YAML content to the temporary file
		err = os.WriteFile(configFile, []byte(yamlData), 0o644)
		Expect(err).ToNot(HaveOccurred())

		// Env variables override
		env := map[string]string{
			"CONFIG_PATH": configFile,
		}
		for name, value := range env {
			Expect(os.Setenv(name, value)).To(Succeed())
		}

		defer func() {
			for name := range env {
				Expect(os.Unsetenv(name)).To(Succeed())
			}
		}()

		// Initialize the config from YAML file
		cfg, err = config.Build()
		Expect(err).To(HaveOccurred())
	})
})

func generateTLSCertificate(certFile, keyFile string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour) // Valid for 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Save certificate to file
	certFileHandle, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certFileHandle.Close()

	err = pem.Encode(certFileHandle, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err != nil {
		return err
	}

	// Save private key to file
	keyFileHandle, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyFileHandle.Close()

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	err = pem.Encode(keyFileHandle, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}

	return nil
}

var _ = Describe("Compression", func() {
	It("should decode valid compression level from string", func() {
		for input, expected := range map[string]config.Compression{
			"none":    config.Compression(gzip.NoCompression),
			"default": config.Compression(gzip.DefaultCompression),
			"speed":   config.Compression(gzip.BestSpeed),
			"best":    config.Compression(gzip.BestCompression),
			"7":       config.Compression(7),
		} {
			output, err := config.DecodeCompression(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		}
	})

	It("should decode valid compression level from number", func() {
		for input, expected := range map[int8]config.Compression{
			0:  config.Compression(gzip.NoCompression),
			-1: config.Compression(gzip.DefaultCompression),
			1:  config.Compression(gzip.BestSpeed),
			9:  config.Compression(gzip.BestCompression),
			7:  config.Compression(7),
		} {
			output, err := config.DecodeCompression(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		}
	})

	It("should return an error for an invalid compression level", func() {
		for _, input := range []any{"invalid", "100", 100, -100} {
			_, err := config.DecodeCompression(reflect.TypeOf(input), reflect.TypeFor[config.Compression](), input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported compression level"))
		}
	})
})

var _ = Describe("Root", func() {
	It("should use existing folder", func() {
		symlink := filepath.Join(tempDir, "symlink")
		Expect(os.Symlink("/tmp", symlink)).To(Succeed())
		for input, expected := range map[string]config.Root{
			"/tmp":        "/tmp",
			"/":           "/",
			"/tmp/../tmp": "/tmp",
			tempDir:       config.Root(tempDir),
			symlink:       config.Root(symlink),
		} {
			output, err := config.DecodeRoot(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		}
	})

	It("should return an error for an invalid root path", func() {
		for _, input := range []any{"invalid", "/etc/passwd", 100, true} {
			_, err := config.DecodeRoot(reflect.TypeOf(input), reflect.TypeFor[config.Root](), input)
			Expect(err).To(HaveOccurred())
		}
	})
})
