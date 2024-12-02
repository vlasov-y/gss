package pkg

// import (
// 	"fmt"
// 	"gss/config"
// 	"os"
// 	"strings"

// 	"github.com/go-acme/lego/v4/certificate"
// 	"github.com/go-acme/lego/v4/lego"
// )

// func IssueACMECertificate(cfg *config.Config) error {
// 	// Initialize Lego ACME client
// 	client, err := lego.NewClient(lego.NewConfig(nil))
// 	if err != nil {
// 		return fmt.Errorf("failed to create ACME client: %w", err)
// 	}

// 	// Set up ACME registration with email
// 	err = client.Register(lego.NewRegistration().SetEmail(cfg.ACMEEmail))
// 	if err != nil {
// 		return fmt.Errorf("ACME registration failed: %w", err)
// 	}

// 	// Obtain domains from the configuration
// 	domains := strings.Split(cfg.ACMEDomains, ",")
// 	certRequest := certificate.ObtainRequest{
// 		Domains: domains,
// 		Bundle:  true,
// 	}

// 	// Obtain the certificate
// 	certificates, err := client.Certificate.Obtain(certRequest)
// 	if err != nil {
// 		return fmt.Errorf("failed to obtain ACME certificate: %w", err)
// 	}

// 	// Save the certificate and private key to files
// 	if err := saveACMECertificate(certificates); err != nil {
// 		return fmt.Errorf("failed to save ACME certificate: %w", err)
// 	}

// 	return nil
// }

// func saveACMECertificate(cert *certificate.Resource) error {
// 	// Save the private key to a file
// 	privKeyPath := "certs/private.key"
// 	privKeyFile, err := os.Create(privKeyPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create private key file: %w", err)
// 	}
// 	defer privKeyFile.Close()

// 	if _, err := privKeyFile.Write(cert.PrivateKey); err != nil {
// 		return fmt.Errorf("failed to write private key: %w", err)
// 	}

// 	// Save the certificate to a file
// 	certPath := "certs/certificate.crt"
// 	certFile, err := os.Create(certPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create certificate file: %w", err)
// 	}
// 	defer certFile.Close()

// 	if _, err := certFile.Write(cert.Certificate); err != nil {
// 		return fmt.Errorf("failed to write certificate: %w", err)
// 	}

// 	// Save the CA certificate to a file
// 	caPath := "certs/ca.crt"
// 	caFile, err := os.Create(caPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create CA certificate file: %w", err)
// 	}
// 	defer caFile.Close()

// 	if _, err := caFile.Write(cert.CA); err != nil {
// 		return fmt.Errorf("failed to write CA certificate: %w", err)
// 	}

// 	return nil
// }
