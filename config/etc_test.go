package config_test

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vlasov-y/gss/config"
)

var _ = Describe("Etc", func() {
	Context("Compression", func() {
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
			for input, expected := range map[any]config.Compression{
				0:         config.Compression(gzip.NoCompression),
				-1:        config.Compression(gzip.DefaultCompression),
				1:         config.Compression(gzip.BestSpeed),
				6:         config.Compression(6),
				uint32(7): config.Compression(7),
				uint64(8): config.Compression(8),
				9:         config.Compression(gzip.BestCompression),
			} {
				output, err := config.DecodeCompression(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should return an error for an invalid compression level", func() {
			for _, input := range []any{
				"invalid",
				"100",
				100,
				-100,
				uint64(0xFFFFFFFFFFFFFFFF),
				true,
			} {
				_, err := config.DecodeCompression(reflect.TypeOf(input), reflect.TypeFor[config.Compression](), input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported compression level"))
			}
		})
	})

	Context("Root", func() {
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

	Context("Headers", func() {
		It("should decode valid YAML/JSON from string", func() {
			for input, expected := range map[string]http.Header{
				"a: b\n": func() http.Header {
					h := http.Header{}
					h.Add("a", "b")
					return h
				}(),
				`{"A":"B:C"}`: func() http.Header {
					h := http.Header{}
					h.Add("A", "B:C")
					return h
				}(),
				"{}":   {},
				"null": {},
			} {
				output, err := config.DecodeHeaders(reflect.TypeOf(input), reflect.TypeOf(expected), input)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(expected))
			}
		})

		It("should return an error for an invalid YAML/JSON", func() {
			for _, input := range []any{
				"a:\n  c: d\n",
				"{",
				"[]",
				true,
				3.14,
				[]string{"test"},
				map[string]any{
					"a":      "b",
					"broken": 3.14,
				},
			} {
				_, err := config.DecodeHeaders(reflect.TypeOf(input), reflect.TypeFor[http.Header](), input)
				Expect(err).To(HaveOccurred())
			}
		})

		It("should return an error for invalid headers", func() {
			for _, input := range []any{
				`{"Invalid !":'B:C'}`,
			} {
				output, err := config.DecodeHeaders(reflect.TypeOf(input), reflect.TypeFor[http.Header](), input)
				fmt.Println(output)
				Expect(err).To(HaveOccurred())
			}
		})
	})
})
