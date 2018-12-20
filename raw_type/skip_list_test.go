package raw_type

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		ret := sl.IsInRange(RangeSpec{
			Min:   0,
			Max:   9.9,
			MinEx: false,
			MaxEx: false,
		})
		Expect(ret).To(Equal(true))

		ret = sl.IsInRange(RangeSpec{
			Min:   10,
			Max:   10.1,
			MinEx: false,
			MaxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(RangeSpec{
			Min:   -2,
			Max:   -1,
			MinEx: false,
			MaxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(RangeSpec{
			Min:   9.8,
			Max:   10.1,
			MinEx: false,
			MaxEx: false,
		})
		Expect(ret).To(Equal(true))

		ret = sl.IsInRange(RangeSpec{
			Min:   -10,
			Max:   0,
			MinEx: false,
			MaxEx: false,
		})
		Expect(ret).To(Equal(true))

		// 开区间
		ret = sl.IsInRange(RangeSpec{
			Min:   9.9,
			Max:   10.1,
			MinEx: true,
			MaxEx: false,
		})
		Expect(ret).To(Equal(false))

		ret = sl.IsInRange(RangeSpec{
			Min:   -1,
			Max:   0.0,
			MinEx: false,
			MaxEx: true,
		})
		Expect(ret).To(Equal(false))
	})

	It("test FirstInRange and LastInRange", func() {
		node := sl.FirstInRange(RangeSpec{
			Min:   1.0,
			Max:   9.9,
			MinEx: false,
			MaxEx: false,
		})
		Expect(node.score).To(Equal(1.0))

		node = sl.FirstInRange(RangeSpec{
			Min:   1.0,
			Max:   9.9,
			MinEx: true,
			MaxEx: false,
		})
		Expect(node.score).To(Equal(1.1))

		node = sl.LastInRange(RangeSpec{
			Min:   1.0,
			Max:   9.9,
			MinEx: false,
			MaxEx: false,
		})
		Expect(node.score).To(Equal(9.9))

		node = sl.LastInRange(RangeSpec{
			Min:   1.0,
			Max:   9.9,
			MinEx: false,
			MaxEx: true,
		})
		Expect(node.score).To(Equal(9.8))
	})

	It("test DeleteRangeByScore", func() {
		count := sl.DeleteRangeByScore(RangeSpec{
			Min:   0.0,
			Max:   0.9,
			MinEx: false,
			MaxEx: false,
		})
		Expect(count).To(Equal(10))

		count = sl.DeleteRangeByScore(RangeSpec{
			Min:   1.0,
			Max:   1.9,
			MinEx: true,
			MaxEx: false,
		})
		Expect(count).To(Equal(9))

		count = sl.DeleteRangeByScore(RangeSpec{
			Min:   2.0,
			Max:   2.9,
			MinEx: false,
			MaxEx: true,
		})
		Expect(count).To(Equal(9))

		count = sl.DeleteRangeByScore(RangeSpec{
			Min:   3.0,
			Max:   3.9,
			MinEx: true,
			MaxEx: true,
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
