package raw_type

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestList(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test List")
}

var _ = Describe("test list ", func() {
	var list *List
	BeforeEach(func() {
		list = ListCreate()
		for i := 9; i >= 0; i-- {
			list = list.ListAddNodeHead(fmt.Sprintf("%d", i))
		}
		Expect(list.length).To(Equal(10))

		for i := 10; i < 20; i++ {
			list = list.ListAddNodeTail(fmt.Sprintf("%d", i))
		}
		Expect(list.length).To(Equal(20))
	})

	It("test iterator from head and tail", func() {
		i := 0
		for node := list.ListFirst(); node != nil; node = node.NodeNext() {
			Expect(node.NodeValue()).To(Equal(fmt.Sprintf("%d", i)))
			i += 1
		}
		i = 19
		for node := list.ListLast(); node != nil; node = node.NodePrev() {
			Expect(node.NodeValue()).To(Equal(fmt.Sprintf("%d", i)))
			i -= 1
		}
	})
	It("test search value in list", func() {
		for i := 0; i < 20; i++ {
			Expect(list.ListSearchKey(fmt.Sprintf("%d", i))).NotTo(Equal(nil))
		}
	})

	It("test ListIndex", func() {
		for i := 0; i < 20; i++ {
			Expect(list.ListIndex(i).NodeValue()).To(Equal(fmt.Sprintf("%d", i)))
		}
		for i := 1; i < 20; i++ {
			Expect(list.ListIndex(-i).NodeValue()).To(Equal(fmt.Sprintf("%d", 20-i)))
		}
	})
})
