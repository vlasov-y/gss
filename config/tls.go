package config

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
)

func DecodeTLSCertificate(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(TLSCertificate{}) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("certificate expects a string, got %T", data)
	}
	input := data.(string)
	// Try to parse input as PEM content directly
	block, _ := pem.Decode([]byte(input))
	if block != nil && block.Type == "CERTIFICATE" {
		parsedCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("invalid certificate content: %w", err)
		}
		return TLSCertificate{Block: block, Cert: parsedCert}, nil
	}
	// If parsing as PEM fails, treat input as a file path
	bytes, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("certificate file not found or unreadable: %w", err)
	}
	block, _ = pem.Decode(bytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid certificate content: file %s does not contain a valid PEM certificate", input)
	}
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate in file %s: %v", input, err)
	}
	return TLSCertificate{Block: block, Cert: parsedCert}, nil
}

func DecodeTLSPrivateKey(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(TLSPrivateKey{}) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("private key expects a string, got %T", data)
	}
	input := data.(string)
	// Try to parse input as PEM content directly
	block, _ := pem.Decode([]byte(input))
	if block != nil {
		privateKey, err := parsePrivateKey(block)
		if err == nil {
			return TLSPrivateKey{Block: block, Key: privateKey}, nil
		}
	}
	// If parsing as PEM fails, treat input as a file path
	bytes, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("private key file not found or unreadable: %w", err)
	}
	block, _ = pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("file %s does not contain valid PEM data", input)
	}
	privateKey, err := parsePrivateKey(block)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return TLSPrivateKey{Block: block, Key: privateKey}, nil
}

// Helper function to parse private keys from a PEM block
func parsePrivateKey(block *pem.Block) (interface{}, error) {
	var privateKey interface{}
	var err error
	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if ed25519Key, ok := privateKey.(ed25519.PrivateKey); ok {
			privateKey = ed25519Key
		} else {
			err = errors.New("unsupported private key type in PKCS#8 format")
		}
	default:
		err = fmt.Errorf("unsupported private key type %q", block.Type)
	}
	return privateKey, err
}

func DecodeTLSVersion(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(TLSVersion(0)) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("TLS version expects a string, got %T", data)
	}
	input := data.(string)
	version := strings.ToUpper(input)
	supported := map[string]TLSVersion{
		"TLS1.0": TLSVersion(tls.VersionTLS10),
		"TLS1.1": TLSVersion(tls.VersionTLS11),
		"TLS1.2": TLSVersion(tls.VersionTLS12),
		"TLS1.3": TLSVersion(tls.VersionTLS13),
	}
	var tlsVersion TLSVersion
	var ok bool
	if tlsVersion, ok = supported[version]; !ok {
		return nil, fmt.Errorf("unsupported TLS version: %s", version)
	}
	return tlsVersion, nil
}

func DecodeTLSCurves(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(TLSCurves{}) {
		return data, nil
	}
	// Function to parse TLS curves
	parseTLSCurves := func(curves []string) (TLSCurves, error) {
		// Map of supported TLS curves
		supported := map[string]tls.CurveID{
			"P-256":   tls.CurveP256,
			"P-384":   tls.CurveP384,
			"P-521":   tls.CurveP521,
			"X25519":  tls.X25519,
			"ED25519": tls.X25519, // ED25519 is also supported as X25519
		}
		var result TLSCurves
		for _, curve := range curves {
			// Trim spaces and convert to uppercase for case-insensitive matching
			curve = strings.ToUpper(strings.TrimSpace(curve))
			if id, ok := supported[curve]; ok {
				if slices.Contains(result, id) {
					return nil, fmt.Errorf("duplicate TLS curves: %s", curve)
				}
				result = append(result, id)
			} else {
				return nil, fmt.Errorf("unsupported TLS curve: %s", curve)
			}
		}
		return result, nil
	}
	// If the input is a string, split it into a slice of curves
	if f.Kind() == reflect.String {
		curveStr := data.(string)
		// Handle comma-separated curves
		curves := strings.Split(curveStr, ",")
		return parseTLSCurves(curves)
	}
	// If the input is already a []string, parse it directly
	if f.Kind() == reflect.Slice {
		curves := []string{}
		for i, object := range data.([]interface{}) {
			s, ok := object.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported type for TLS curves at index %d, expected string, got %s", i, reflect.TypeOf(object))
			}
			curves = append(curves, s)
		}
		return parseTLSCurves(curves)
	}
	// If not a string or []string, return an error
	return nil, fmt.Errorf("unsupported type for TLS curves, expected string or []string, got %s", f.Kind())
}

func DecodeTLSCiphers(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	// Check if the target type is TLSCiphers
	if t != reflect.TypeOf(TLSCiphers{}) {
		return data, nil
	}
	// Function to parse TLS ciphers
	parseTLSCiphers := func(ciphers []string) (TLSCiphers, error) {
		// Map of supported TLS ciphers
		supported := map[string]uint16{
			"TLS_AES_128_GCM_SHA256":                        tls.TLS_AES_128_GCM_SHA256,
			"TLS_AES_256_GCM_SHA384":                        tls.TLS_AES_256_GCM_SHA384,
			"TLS_CHACHA20_POLY1305_SHA256":                  tls.TLS_CHACHA20_POLY1305_SHA256,
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			"TLS_RSA_WITH_AES_128_CBC_SHA256":               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			"TLS_RSA_WITH_AES_128_GCM_SHA256":               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			"TLS_RSA_WITH_AES_256_GCM_SHA384":               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		}
		var result TLSCiphers
		for _, c := range ciphers {
			if cipher, ok := supported[c]; ok {
				if slices.Contains(result, cipher) {
					return nil, fmt.Errorf("duplicate TLS cipher suite: %s", c)
				}
				result = append(result, cipher)
			} else {
				return nil, fmt.Errorf("unsupported TLS cipher suite: %s", c)
			}
		}
		return result, nil
	}
	// If the input is a string, split it into a slice of ciphers
	if f.Kind() == reflect.String {
		cipherStr := data.(string)
		// Handle comma-separated ciphers
		ciphers := strings.Split(cipherStr, ",")
		return parseTLSCiphers(ciphers)
	}
	// If the input is already a []string, parse it directly
	if f.Kind() == reflect.Slice {
		ciphers := []string{}
		for i, object := range data.([]interface{}) {
			s, ok := object.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported type for TLS ciphers at index %d, expected string, got %s", i, reflect.TypeOf(object))
			}
			ciphers = append(ciphers, s)
		}
		return parseTLSCiphers(ciphers)
	}
	// If not a string or []string, return an error
	return nil, fmt.Errorf("unsupported type for TLS ciphers, expected string or []string, got %s", f.Kind())
}
