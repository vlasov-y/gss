package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/vlasov-y/gss/config"
	"github.com/vlasov-y/gss/pkg"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	. "github.com/onsi/gomega"
)

func init() {
	RegisterFailHandler(func(message string, _ ...int) {
		log.Fatalln(message)
	})
}

func main() {
	// Load configuration
	cfg, err := config.Build()
	Expect(pkg.PrefixError("failed to build a configuration object", err)).NotTo(HaveOccurred())

	// Create router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files safely
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFile(w, r, string(cfg.Root))
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Setup TLS if required
	if cfg.TLS.Certificate.Block != nil && cfg.TLS.Key.Block != nil {
		server.TLSConfig, err = setupTLSConfig(cfg.TLS)
		if err != nil {
			log.Fatalf("failed to configure TLS: %v", err)
		}
	}

	// Start server in a goroutine
	go func() {
		if server.TLSConfig != nil {
			log.Printf("Starting HTTPS server on port %d", cfg.Port)
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server failed: %v", err)
			}
		} else {
			log.Printf("Starting HTTP server on port %d", cfg.Port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server failed: %v", err)
			}
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

// setupTLSConfig sets up the TLS configuration based on the provided settings.
func setupTLSConfig(tlsConfig config.TLSConfig) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(pem.EncodeToMemory(tlsConfig.Certificate.Block), pem.EncodeToMemory(tlsConfig.Key.Block))
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS key pair: %w", err)
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Enable client certificate authentication if CA is set
	if tlsConfig.CA.Block != nil {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pem.EncodeToMemory(tlsConfig.CA.Block)) {
			return nil, errors.New("failed to parse CA certificate")
		}
		tlsConf.ClientCAs = certPool
		tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConf, nil
}

// gracefulShutdown handles server shutdown on stop signals.
func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func serveStaticFile(w http.ResponseWriter, r *http.Request, rootDir string) {
	// Ensure the cleanPath is within the rootDir to prevent directory traversal
	fullPath := filepath.Join(rootDir, filepath.Clean(r.URL.Path))
	if fullPath != rootDir && !strings.HasPrefix(fullPath, rootDir+string(os.PathSeparator)) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Check if the path is a directory
	info, err := os.Lstat(fullPath)
	if err != nil || !info.Mode().IsRegular() {
		http.ServeFile(w, r, filepath.Join(rootDir, "index.html"))
		return
	}

	// If it's a file, serve it
	http.ServeFile(w, r, fullPath)
}
