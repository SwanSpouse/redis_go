package raw_type

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("test dictEntry equal", func() {
	It("should equal", func() {
		entry1 := NewDictEntry("Key", "Value", nil)
		entry2 := NewDictEntry("Key", "Value", nil)
		Expect(entry1.equals(entry2)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry2.hasCode()))

		entry3 := NewDictEntry("Key", "Value", entry1)
		Expect(entry1.equals(entry3)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry3.hasCode()))

		var entry4 dictEntry
		entry4.Key = "Key"
		entry4.Value = "Value"
		Expect(entry1.equals(entry4)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry4.hasCode()))
	})

	It("should not equal", func() {
		entry1 := NewDictEntry("key1", "Value", nil)
		entry2 := NewDictEntry("key2", "Value", nil)
		Expect(entry1.equals(entry2)).To(Equal(false))
		Expect(entry1.hasCode()).NotTo(Equal(entry2.hasCode()))

		entry3 := NewDictEntry("Key", "Value", nil)
		var entry4 dictEntry
		Expect(entry3.equals(entry4)).To(Equal(false))
	})

	It("compare diff type compare", func() {
		entry1 := NewDictEntry("Key", "Value", nil)
		entry2 := map[string]string{
			"Key": "Key", "Value": "Value",
		}
		Expect(entry1.equals(entry2)).To(Equal(false))
	})

})

func TestDict(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Dict")
}
