package config_test

import (
	"crypto/tls"
	"gss/config"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TLS parsing", func() {
	Context("TLSCertificate", func() {
		var cert config.TLSCertificate

		BeforeEach(func() {
			cert = config.TLSCertificate{}
		})

		It("should decode a valid certificate from file", func() {
			output, err := config.DecodeTLSCertificate(reflect.TypeOf(certFile), reflect.TypeOf(cert), certFile)
			Expect(err).ToNot(HaveOccurred())

			cert, ok = output.(config.TLSCertificate)
			Expect(ok).To(BeTrue())
		})

		It("should decode a valid certificate from text", func() {
			output, err := config.DecodeTLSCertificate(reflect.TypeOf(certContent), reflect.TypeOf(cert), certContent)
			Expect(err).ToNot(HaveOccurred())

			cert, ok = output.(config.TLSCertificate)
			Expect(ok).To(BeTrue())
		})

		It("should return an error for a non-existing certificate file", func() {
			_, err := config.DecodeTLSCertificate(reflect.TypeFor[string](), reflect.TypeOf(cert), "doesNotExist")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("certificate file not found"))
		})

		It("should return an error for a wrong input type", func() {
			_, err := config.DecodeTLSCertificate(reflect.TypeFor[int](), reflect.TypeOf(cert), 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("certificate expects a string"))
		})

		It("should return an error for an invalid certificate data", func() {
			_, err := config.DecodeTLSCertificate(reflect.TypeOf(keyFile), reflect.TypeOf(cert), keyFile)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid certificate content"))
		})
	})

	Context("TLSPrivateKey", func() {
		var key config.TLSPrivateKey

		BeforeEach(func() {
			key = config.TLSPrivateKey{}
		})

		It("should decode a valid private key from file", func() {
			output, err := config.DecodeTLSPrivateKey(reflect.TypeOf(keyFile), reflect.TypeOf(key), keyFile)
			Expect(err).ToNot(HaveOccurred())

			key, ok = output.(config.TLSPrivateKey)
			Expect(ok).To(BeTrue())
		})

		It("should decode a valid private key from text", func() {
			output, err := config.DecodeTLSPrivateKey(reflect.TypeOf(keyContent), reflect.TypeOf(key), keyContent)
			Expect(err).ToNot(HaveOccurred())

			key, ok = output.(config.TLSPrivateKey)
			Expect(ok).To(BeTrue())
		})

		It("should return an error for a non-existing certificate file", func() {
			_, err := config.DecodeTLSPrivateKey(reflect.TypeFor[string](), reflect.TypeOf(key), "doesNotExist")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("private key file not found"))
		})

		It("should return an error for a wrong input type", func() {
			_, err := config.DecodeTLSPrivateKey(reflect.TypeFor[int](), reflect.TypeOf(key), 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("private key expects a string"))
		})

		It("should return an error for an invalid private key data", func() {
			_, err := config.DecodeTLSPrivateKey(reflect.TypeOf(certFile), reflect.TypeOf(key), certFile)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to parse private key"))
		})
	})

	Context("TLSVersion", func() {
		It("should decode valid TLS versions", func() {
			for input, expected := range map[string]config.TLSVersion{
				"TLS1.0": config.TLSVersion(tls.VersionTLS10),
				"TLS1.1": config.TLSVersion(tls.VersionTLS11),
				"TLS1.2": config.TLSVersion(tls.VersionTLS12),
				"TLS1.3": config.TLSVersion(tls.VersionTLS13),
			} {
				output, err := config.DecodeTLSVersion(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should return an error for an unsupported TLS version", func() {
			for _, input := range []string{"SSLv3", "TLS1.5"} {
				_, err := config.DecodeTLSVersion(reflect.TypeOf(input), reflect.TypeFor[config.TLSVersion](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported TLS version"))
			}
		})
	})

	Context("TLSCurves", func() {
		It("should decode valid TLS curves from single string value", func() {
			for input, expected := range map[string]config.TLSCurves{
				"P-256":   {tls.CurveP256},
				"P-384":   {tls.CurveP384},
				"P-521":   {tls.CurveP521},
				"X25519":  {tls.X25519},
				"ED25519": {tls.X25519},
			} {
				output, err := config.DecodeTLSCurves(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid TLS curves from CSV", func() {
			for input, expected := range map[string]config.TLSCurves{
				"P-256,P-384,P-521,X25519": {
					tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519,
				},
			} {
				output, err := config.DecodeTLSCurves(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid TLS curves from []string", func() {
			input := []any{"P-256", "P-384", "P-521", "X25519"}
			expected := config.TLSCurves{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519}
			output, err := config.DecodeTLSCurves(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		})

		It("should return an error for an unsupported TLS curve", func() {
			for _, input := range []string{
				"P-XXX", ",",
			} {
				_, err := config.DecodeTLSCurves(reflect.TypeOf(input), reflect.TypeFor[config.TLSCurves](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported TLS curve"))
			}
		})

		It("should return an error for duplicate TLS curves", func() {
			input := "P-256,P-384,P-256"
			_, err := config.DecodeTLSCurves(reflect.TypeOf(input), reflect.TypeFor[config.TLSCurves](), input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate TLS curves"))
		})
	})

	Context("TLSCipherSuites", func() {
		It("should decode valid TLS cipher suites from single string value", func() {
			for input, expected := range map[string]config.TLSCiphers{
				"TLS_AES_128_GCM_SHA256":                        {tls.TLS_AES_128_GCM_SHA256},
				"TLS_AES_256_GCM_SHA384":                        {tls.TLS_AES_256_GCM_SHA384},
				"TLS_CHACHA20_POLY1305_SHA256":                  {tls.TLS_CHACHA20_POLY1305_SHA256},
				"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          {tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA},
				"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       {tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256},
				"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       {tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
				"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          {tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA},
				"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       {tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
				"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": {tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256},
				"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              {tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA},
				"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           {tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA},
				"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            {tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
				"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         {tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256},
				"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         {tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
				"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            {tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
				"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         {tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
				"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   {tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256},
				"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                {tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA},
				"TLS_RSA_WITH_AES_128_CBC_SHA256":               {tls.TLS_RSA_WITH_AES_128_CBC_SHA256},
				"TLS_RSA_WITH_AES_128_GCM_SHA256":               {tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
				"TLS_RSA_WITH_AES_256_GCM_SHA384":               {tls.TLS_RSA_WITH_AES_256_GCM_SHA384},
			} {
				output, err := config.DecodeTLSCiphers(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid TLS cipher suites from CSV", func() {
			for input, expected := range map[string]config.TLSCiphers{
				"TLS_AES_128_GCM_SHA256,TLS_AES_256_GCM_SHA384": {
					tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384,
				},
			} {
				output, err := config.DecodeTLSCiphers(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid TLS cipher suites from []string", func() {
			input := []any{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384"}
			expected := config.TLSCiphers{tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384}
			output, err := config.DecodeTLSCiphers(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		})

		It("should return an error for an unsupported TLS cipher suite or wrong input", func() {
			for _, input := range []string{
				"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
				"TLS_AES_128_GCM_SHA256 ",
				"tls_AES_128_GCM_SHA256",
				"incorrect",
				",",
			} {
				_, err := config.DecodeTLSCiphers(reflect.TypeOf(input), reflect.TypeFor[config.TLSCiphers](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported TLS cipher suite"))
			}
		})

		It("should return an error for duplicate TLS cipher suites", func() {
			input := "TLS_AES_128_GCM_SHA256,TLS_AES_128_GCM_SHA256"
			_, err := config.DecodeTLSCiphers(reflect.TypeOf(input), reflect.TypeFor[config.TLSCiphers](), input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate TLS cipher suite"))
		})
	})
})
