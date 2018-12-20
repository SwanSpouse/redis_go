package util

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNumberUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Number Util")
}

var _ = Describe("FormatFloatString", func() {
	It("FormatFloatString", func() {
		ret, err := FormatFloatString("")
		Expect(err).NotTo(Equal(nil))
		Expect(ret).To(Equal(""))

		ret, _ = FormatFloatString("100")
		Expect(ret).To(Equal("100"))

		ret, _ = FormatFloatString("5.0")
		Expect(ret).To(Equal("5.00"))

		ret, _ = FormatFloatString("5.00")
		Expect(ret).To(Equal("5.00"))

		ret, _ = FormatFloatString("5.00000")
		Expect(ret).To(Equal("5.00"))

		ret, _ = FormatFloatString("5.0123410000")
		Expect(ret).To(Equal("5.012341"))

		ret, _ = FormatFloatString("5.10")
		Expect(ret).To(Equal("5.10"))

		ret, _ = FormatFloatString("5.10000")
		Expect(ret).To(Equal("5.10"))
	})
})
