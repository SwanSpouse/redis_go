package raw_type

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
)

func TestNewSkipList(t *testing.T) {
	skipList := NewSkipList()

	for _, item := range skipList.header.level {
		fmt.Printf("forward:%+v\n", item.forward)
		fmt.Printf("span:%+v\n", item.span)
	}
}

func TestSkipList(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test SkipList")
}

var _ = Describe("test dictEntry equal", func() {
	It("should equal", func() {
		entry1 := NewDictEntry(0, "Key", "Value", nil)
		entry2 := NewDictEntry(0, "Key", "Value", nil)
		Expect(entry1.equals(entry2)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry2.hasCode()))

		entry3 := NewDictEntry(0, "Key", "Value", entry1)
		Expect(entry1.equals(entry3)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry3.hasCode()))

		var entry4 dictEntry
		entry4.Key = "Key"
		entry4.Value = "Value"
		Expect(entry1.equals(entry4)).To(Equal(true))
		Expect(entry1.hasCode()).To(Equal(entry4.hasCode()))
	})
})
