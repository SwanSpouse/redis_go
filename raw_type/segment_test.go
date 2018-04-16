package raw_type

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"testing"
)

func TestSegment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Dict")
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

	It("should not equal", func() {
		entry1 := NewDictEntry(0, "key1", "Value", nil)
		entry2 := NewDictEntry(0, "key2", "Value", nil)
		Expect(entry1.equals(entry2)).To(Equal(false))
		Expect(entry1.hasCode()).NotTo(Equal(entry2.hasCode()))

		entry3 := NewDictEntry(0, "Key", "Value", nil)
		var entry4 dictEntry
		Expect(entry3.equals(entry4)).To(Equal(false))
	})

	It("compare diff type compare", func() {
		entry1 := NewDictEntry(0, "Key", "Value", nil)
		entry2 := map[string]string{
			"Key": "Key", "Value": "Value",
		}
		Expect(entry1.equals(entry2)).To(Equal(false))
	})

})

var _ = Describe("test segment", func() {
	It("test segment printForDebug", func() {
		seg := newSegment(20, 0.75, 300)
		lastModCount := seg.modCount
		for i := 1; i <= 30; i++ {
			seg.put(hash(i), fmt.Sprintf("key%d", i), fmt.Sprintf("%d", i))
		}
		Expect(seg.count).To(Equal(30))
		Expect(seg.modCount).NotTo(Equal(lastModCount))
	})

	It("test segment for concurrent put", func() {
		var wg sync.WaitGroup
		seg := newSegment(10, 0.75, 400)
		wg.Add(5)
		for threadCount := 1; threadCount <= 5; threadCount++ {
			go func(no int) {
				defer wg.Done()
				for i := 1; i <= 30; i++ {
					seg.put(hash(i), fmt.Sprintf("%d-%d", no, i), "")
				}
			}(threadCount)
		}
		wg.Wait()
		Expect(seg.count).To(Equal(150))
	})

})

var _ = Describe("test for remove", func() {
	seg := newSegment(400, 0.75, 300)

	BeforeEach(func() {
		var wg sync.WaitGroup
		wg.Add(5)
		for no := 1; no <= 5; no++ {
			go func(threadNo int) {
				defer wg.Done()
				for i := 0; i < 30; i++ {
					seg.put(hash(i), fmt.Sprintf("%d-%d", threadNo, i+1), "")
				}
			}(no)
		}
		wg.Wait()
		Expect(seg.count).To(Equal(150))
		//seg.printSegForDebug()
	}, 0)

	It("test segment for remove", func() {
		for no := 1; no <= 5; no++ {
			for i := 0; i < 30; i++ {
				seg.remove(hash(i), fmt.Sprintf("%d-%d", no, i+1), "")
			}
			//seg.printSegForDebug()
			Expect(seg.count).To(Equal(150 - 30*no))
		}
	})

	It("test segment for concurrent remove", func() {
		var wg sync.WaitGroup
		wg.Add(5)
		for no := 1; no <= 5; no++ {
			go func(threadNo int) {
				defer wg.Done()
				for i := 0; i < 30; i++ {
					seg.remove(hash(i), fmt.Sprintf("%d-%d", threadNo, i+1), "")
				}
			}(no)
		}
		wg.Wait()
		Expect(seg.count).To(Equal(0))
	})

	It("test segment for replace", func() {
		for no := 1; no <= 5; no++ {
			for i := 0; i < 30; i++ {
				seg.replace(hash(i), fmt.Sprintf("%d-%d", no, i+1), "", "HA")
			}
		}
		Expect(seg.count).To(Equal(150))
		for i := 0; i < len(seg.table); i++ {
			if seg.table[i] == nil {
				continue
			}
			for e := seg.table[i]; e != nil; e = e.next {
				Expect(seg.table[i].Value).To(Equal("HA"))
			}
		}
	})

	It("test segment for concurrent replace", func() {
		var wg sync.WaitGroup
		wg.Add(5)
		for no := 1; no <= 5; no++ {
			go func(threadNo int) {
				defer wg.Done()
				for i := 0; i < 30; i++ {
					seg.replace(hash(i), fmt.Sprintf("%d-%d", threadNo, i+1), "", "HA")
				}
			}(no)
		}
		wg.Wait()
		for i := 0; i < len(seg.table); i++ {
			if seg.table[i] == nil {
				continue
			}
			for e := seg.table[i]; e != nil; e = e.next {
				Expect(seg.table[i].Value).To(Equal("HA"))
			}
		}
	})

	It("test segment for Clear", func() {
		seg.clear()
		Expect(seg.count).To(Equal(0))
		for i := 0; i < len(seg.table); i++ {
			Expect(seg.table[i]).To(BeNil())
		}
	})

	It("test segment for concurrent Clear", func() {
		var wg sync.WaitGroup
		wg.Add(5)
		for no := 1; no <= 5; no++ {
			go func() {
				defer wg.Done()
				seg.clear()
			}()
		}
		wg.Wait()
		Expect(seg.count).To(Equal(0))
		for i := 0; i < len(seg.table); i++ {
			Expect(seg.table[i]).To(BeNil())
		}
	})
})
