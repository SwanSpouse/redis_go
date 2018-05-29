package raw_type

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDebugSkipList(t *testing.T) {
	sl := NewSkipList()

	for i := 0; i < 10; i++ {
		sl.Insert(fmt.Sprintf("x%d", i), float64(i)/10)
	}
	sl.DebugSkipListInfo()
}

func TestSkipList(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test SkipList")
}

var _ = Describe("test skip list", func() {
	var sl *SkipList
	inputSize := 100
	BeforeEach(func() {
		sl = NewSkipList()
		for i := 0; i < inputSize; i++ {
			Expect(sl.length).To(Equal(i))
			sl.Insert(fmt.Sprintf("%d", i), float64(i)/10)
		}
		Expect(sl.length).To(Equal(inputSize))
	})

	AfterEach(func() {
		sl = nil
	})

	It("test delete", func() {
		Expect(sl.length).To(Equal(inputSize))

		for i := inputSize - 1; i >= 0; i-- {
			sl.Delete(fmt.Sprintf("%d", i), float64(i)/10)
			Expect(sl.length).To(Equal(i))
		}
	})

	It("test IsInRange", func() {
		// 闭区间
		ret := sl.IsInRange(rangeSpec{
			min:   0,
			max:   9.9,
			minEx: false,
			maxEx: false,
		})
		Expect(ret).To(Equal(true))

		ret = sl.IsInRange(rangeSpec{
			min:   10,
			max:   10.1,
			minEx: false,
			maxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(rangeSpec{
			min:   -2,
			max:   -1,
			minEx: false,
			maxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(rangeSpec{
			min:   9.8,
			max:   10.1,
			minEx: false,
			maxEx: false,
		})
		Expect(ret).To(Equal(true))

		ret = sl.IsInRange(rangeSpec{
			min:   -10,
			max:   0,
			minEx: false,
			maxEx: false,
		})
		Expect(ret).To(Equal(true))

		// 开区间
		ret = sl.IsInRange(rangeSpec{
			min:   9.9,
			max:   10.1,
			minEx: true,
			maxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(rangeSpec{
			min:   -1,
			max:   0.0,
			minEx: false,
			maxEx: true,
		})
		Expect(ret).To(Equal(false))
	})

	It("test FirstInRange and LastInRange", func() {
		node := sl.FirstInRange(rangeSpec{
			min:   1.0,
			max:   9.9,
			minEx: false,
			maxEx: false,
		})
		Expect(node.score).To(Equal(1.0))

		node = sl.FirstInRange(rangeSpec{
			min:   1.0,
			max:   9.9,
			minEx: true,
			maxEx: false,
		})
		Expect(node.score).To(Equal(1.1))

		node = sl.LastInRange(rangeSpec{
			min:   1.0,
			max:   9.9,
			minEx: false,
			maxEx: false,
		})
		Expect(node.score).To(Equal(9.9))

		node = sl.LastInRange(rangeSpec{
			min:   1.0,
			max:   9.9,
			minEx: false,
			maxEx: true,
		})
		Expect(node.score).To(Equal(9.8))
	})

	It("test DeleteRangeByScore", func() {
		count := sl.DeleteRangeByScore(rangeSpec{
			min:   0.0,
			max:   0.9,
			minEx: false,
			maxEx: false,
		})
		Expect(count).To(Equal(10))

		count = sl.DeleteRangeByScore(rangeSpec{
			min:   1.0,
			max:   1.9,
			minEx: true,
			maxEx: false,
		})
		Expect(count).To(Equal(9))

		count = sl.DeleteRangeByScore(rangeSpec{
			min:   2.0,
			max:   2.9,
			minEx: false,
			maxEx: true,
		})
		Expect(count).To(Equal(9))

		count = sl.DeleteRangeByScore(rangeSpec{
			min:   3.0,
			max:   3.9,
			minEx: true,
			maxEx: true,
		})
		Expect(count).To(Equal(8))
	})

	It("test DeleteRangeByRank", func() {
		count := sl.DeleteRangeByRank(1, 10)
		Expect(count).To(Equal(10))
	})

	It("test GetRank", func() {
		for i := 0; i < inputSize; i++ {
			rank := sl.GetRank(fmt.Sprintf("%d", i), float64(i)/10)
			Expect(i + 1).To(Equal(rank))
		}
	})

	It("test GetElementByRank", func() {
		for i := 0; i < inputSize; i++ {
			node := sl.GetElementByRank(i + 1)
			Expect(fmt.Sprintf("%d", i)).To(Equal(node.obj))
		}
	})
})
