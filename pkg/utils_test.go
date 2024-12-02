package pkg_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vlasov-y/gss/pkg"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("PrefixError", func() {
	It("should return an error with prefixed text", func() {
		err := pkg.PrefixError("prefix", errors.New("test"))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(HavePrefix("prefix: "))
	})
})
