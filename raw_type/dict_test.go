package raw_type

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestWriteTest(t *testing.T) {
	// 编写Test的时候在这里写，写好了再迁移到Describe里
}

func TestDict(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Dict")
}

var _ = Describe("test dict base operations", func() {
	dict := NewDict()
	inputSize := 100
	BeforeEach(func() {
		for i := 0; i < inputSize; i++ {
			dict.Put(i, i)
		}
		Expect(dict.Size()).To(Equal(inputSize))
	})

	It("test dict operation get", func() {
		for i := 0; i < inputSize; i++ {
			value := dict.Get(i)
			Expect(value).To(Equal(i))
		}
		Expect(dict.Size()).To(Equal(inputSize))
	})

	It("test dict operation IsEmpty", func() {
		Expect(dict.IsEmpty()).To(BeFalse())
		dict.Clear()
		Expect(dict.IsEmpty()).To(BeTrue())
	})

	It("test dict operation replace", func() {
		for i := 0; i < inputSize; i++ {
			isSucc := dict.Replace(i, i, i+1)
			Expect(isSucc).To(BeTrue())
		}
		for i := 0; i < inputSize; i++ {
			newValue := dict.Get(i)
			Expect(newValue).To(Equal(i + 1))
		}
	})

	It("test dict operation remove", func() {
		for i := 0; i < inputSize; i++ {
			oldValue := dict.Remove(i, i)
			Expect(oldValue).To(Equal(i))
		}
		for i := 0; i < inputSize; i++ {
			value := dict.Get(i)
			Expect(value).To(BeNil())
		}
	})

	It("test dict operation Clear", func() {
		Expect(dict.Size()).To(Equal(inputSize))
		dict.Clear()
		Expect(dict.Size()).To(Equal(0))
	})

	It("test dict operation ContainsKey", func() {
		for i := 0; i < inputSize; i++ {
			Expect(dict.ContainsKey(i)).To(BeTrue())
		}
		dict.Clear()
		for i := 0; i < inputSize; i++ {
			Expect(dict.ContainsKey(i)).To(BeFalse())
		}
	})

	It("test dict operation ContainsValue || Contains", func() {
		for i := 0; i < inputSize; i++ {
			Expect(dict.ContainsValue(i)).To(BeTrue())
		}
		dict.Clear()
		for i := 0; i < inputSize; i++ {
			Expect(dict.ContainsValue(i)).To(BeFalse())
		}
	})

})
