package config_test

import (
	"reflect"

	"github.com/vlasov-y/gss/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ACME", func() {
	Context("ACMEEmail", func() {
		It("should decode a valid email", func() {
			input := "valid@example.com"
			output, err := config.DecodeACMEEmail(reflect.TypeFor[string](), reflect.TypeFor[config.ACMEEmail](), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(config.ACMEEmail(input)))
		})

		It("should return an error for an invalid email", func() {
			for _, input := range []any{
				"invalid", true, 3.14, nil, []string{"invalid"},
			} {
				_, err := config.DecodeACMEEmail(reflect.TypeOf(input), reflect.TypeFor[config.ACMEEmail](), input)
				Expect(err).To(HaveOccurred())
			}
		})
	})

	Context("ACMEURL", func() {
		It("should decode a valid URL", func() {
			input := "https://example.com"
			output, err := config.DecodeACMEURL(reflect.TypeFor[string](), reflect.TypeFor[config.ACMEURL](), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(config.ACMEURL(input)))
		})

		It("should return an error for an invalid URL", func() {
			for _, input := range []any{"invalid", true, 3.14} {
				_, err := config.DecodeACMEURL(reflect.TypeOf(input), reflect.TypeFor[config.ACMEURL](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid ACME URL"))
			}
		})
	})

	Context("ACMEDomains", func() {
		It("should decode valid domains from single string value", func() {
			for input, expected := range map[string]config.ACMEDomains{
				"example.com":   {"example.com"},
				"*.example.com": {"*.example.com"},
			} {
				output, err := config.DecodeACMEDomains(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid domains from CSV", func() {
			for input, expected := range map[string]config.ACMEDomains{
				"example.com,*.example.com": {
					"example.com", "*.example.com",
				},
			} {
				output, err := config.DecodeACMEDomains(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should decode valid domains from []string", func() {
			input := []any{"example.com", "*.example.com"}
			expected := config.ACMEDomains{"example.com", "*.example.com"}
			output, err := config.DecodeACMEDomains(reflect.TypeOf(input), reflect.TypeOf(expected), input)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		})

		It("should return an error for invalid domain", func() {
			for _, input := range []any{
				"com", ",", "", "*", "*.*", "invalid..com",
				[]int{1, 2, 3},
				true, 3.14,
			} {
				_, err := config.DecodeACMEDomains(reflect.TypeOf(input), reflect.TypeFor[config.ACMEDomains](), input)
				Expect(err).To(HaveOccurred())
			}
		})

		It("should return an error for duplicate domains", func() {
			input := "example.com,example.com"
			_, err := config.DecodeACMEDomains(reflect.TypeOf(input), reflect.TypeFor[config.ACMEDomains](), input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate domain name"))
		})
	})

	Context("ACMEChallengePath", func() {
		It("should decode a valid challenge path", func() {
			for _, input := range []string{
				"/", "/test", "/test/", "/test/test",
			} {
				output, err := config.DecodeACMEChallengePath(reflect.TypeOf(input), reflect.TypeFor[config.ACMEChallengePath](), input)
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(Equal(config.ACMEChallengePath(input)))
			}
		})

		It("should return an error for an invalid challenge path", func() {
			for _, input := range []any{
				"", "test", "test/", "/!@#$%", true, 3.14,
			} {
				_, err := config.DecodeACMEChallengePath(reflect.TypeOf(input), reflect.TypeFor[config.ACMEChallengePath](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid ACME challenge path"))
			}
		})
	})
})
