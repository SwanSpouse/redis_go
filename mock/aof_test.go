package mock

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
)

var _ = Describe("Redis aof test", func() {
	It("generate data to aof file", func() {
		if !LoadDataFromAofFile {
			return
		}
		cn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())

		w := NewRequestWriter(cn)
		r := NewResponseReader(cn)

		//  generate string data
		input := make([][]byte, 0)
		for i := 0; i < 100; i++ {
			input = append(input, []byte(fmt.Sprintf("aof#generate-string-key%d", i)))
			input = append(input, []byte(fmt.Sprintf("value%d", i)))
		}
		w.WriteCmd(handlers.RedisStringCommandMSet, input...)
		w.Flush()

		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		// generate list data
		listInput := make([]string, 0)
		listInput = append(listInput, "aof#generate-list")
		for i := 0; i < 100; i++ {
			listInput = append(listInput, fmt.Sprintf("value%d", i))
		}
		w.WriteCmdString(handlers.RedisListCommandRPush, listInput...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("100"))

		// generate hash data
		hashInput := make([]string, 0)
		hashInput = append(hashInput, "aof#generate-hash")
		for i := 0; i < 100; i++ {
			hashInput = append(hashInput, fmt.Sprintf("key%d", i))
			hashInput = append(hashInput, fmt.Sprintf("value%d", i))
		}
		w.WriteCmdString(handlers.RedisHashCommandHMSet, hashInput...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		// generate set data
		setInput := make([]string, 0)
		setInput = append(setInput, "aof#generate-set")
		for i := 0; i < 100; i++ {
			setInput = append(setInput, fmt.Sprintf("item%d", i))
		}
		w.WriteCmdString(handlers.RedisSetCommandSADD, setInput...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("100"))

		// generate sorted set data
		sortedSetInput := make([]string, 0)
		sortedSetInput = append(sortedSetInput, "aof#generate-sorted-set")
		for i := 0; i < 100; i++ {
			sortedSetInput = append(sortedSetInput, fmt.Sprintf("%.2f", float64(i)/100))
			sortedSetInput = append(sortedSetInput, fmt.Sprintf("item%d", i))
		}
		w.WriteCmdString(handlers.RedisSortedSetCommandZAdd, sortedSetInput...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("100"))
	})
})
