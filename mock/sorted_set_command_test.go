package mock

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
	"redis_go/server"
	"redis_go/util"
)

var _ = Describe("TestRedisSortedSetCommands", func() {
	var w *RequestWriter
	var r *ResponseReader

	sortedSetCommandTestBaseKey := "sorted_set_command_test_base_key"
	BeforeEach(func() {
		cn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())

		w = NewRequestWriter(cn)
		r = NewResponseReader(cn)

		// first truncate all DB
		w.WriteCmdString(server.RedisServerCommandFlushAll)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		// init basic set
		input := make([]string, 0)
		input = append(input, sortedSetCommandTestBaseKey)
		for i := 0; i < 10; i++ {
			input = append(input, fmt.Sprintf("%.2f", float64(i)/100))
			input = append(input, fmt.Sprintf("item%d", i))
		}
		w.WriteCmdString(handlers.RedisSortedSetCommandZAdd, input...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))
	})

	AfterEach(func() {
		w.WriteCmdString(handlers.RedisKeyCommandDel, sortedSetCommandTestBaseKey)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))
	})

	It("Test redis sorted set command ZAdd ZCard", func() {
		w.WriteCmdString(handlers.RedisSortedSetCommandZCard, sortedSetCommandTestBaseKey)
		w.Flush()

		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))

		input := make([]string, 0)
		input = append(input, sortedSetCommandTestBaseKey)
		for i := 10; i < 20; i++ {
			input = append(input, fmt.Sprintf("%.2f", float64(i)/100))
			input = append(input, fmt.Sprintf("item%d", i))
		}
		w.WriteCmdString(handlers.RedisSortedSetCommandZAdd, input...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))

		w.WriteCmdString(handlers.RedisSortedSetCommandZCard, sortedSetCommandTestBaseKey)
		w.Flush()

		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("20"))
	})

	It("Test redis sorted set command ZCount", func() {
		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisSortedSetCommandZCount, sortedSetCommandTestBaseKey, "0", fmt.Sprintf("%.2f", float64(i)/100))
			w.Flush()

			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("%d", i+1)))
		}
		for i := 9; i >= 0; i-- {
			w.WriteCmdString(handlers.RedisSortedSetCommandZCount, sortedSetCommandTestBaseKey, fmt.Sprintf("%.2f", float64(i)/100), "0.09")
			w.Flush()

			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("%d", 10-i)))
		}
	})

	It("Test redis sorted set command ZIncrBy", func() {
		for i := 1; i < 5; i++ {
			w.WriteCmdString(handlers.RedisSortedSetCommandZIncrBy, sortedSetCommandTestBaseKey, fmt.Sprintf("item%d", i), fmt.Sprintf("%.2f", float64(i)/100))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("%.2f", float64(i*2)/100)))
		}
		for i := 1; i < 5; i++ {
			w.WriteCmdString(handlers.RedisSortedSetCommandZIncrBy, sortedSetCommandTestBaseKey, fmt.Sprintf("new_item%d", i), fmt.Sprintf("%.2f", float64(i)/100))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("%.2f", float64(i)/100)))
		}
	})

	It("Test redis sorted set command ZRem", func() {
		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisSortedSetCommandZRem, sortedSetCommandTestBaseKey, fmt.Sprintf("item%d", i), fmt.Sprintf("new_item%d", i))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("1"))
		}
		w.WriteCmdString(handlers.RedisSortedSetCommandZCard, sortedSetCommandTestBaseKey)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))
	})

	It("Test redis sorted set command ZRange", func() {
		for endIndex := 0; endIndex < 10; endIndex++ {
			w.WriteCmdString(handlers.RedisSortedSetCommandZRange, sortedSetCommandTestBaseKey, "0", fmt.Sprintf("%d", endIndex))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())

			index := 0
			for i := 0; i <= endIndex; i++ {
				Expect(ret[index]).To(Equal(fmt.Sprintf("item%d", i)))
				index += 1
				Expect(ret[index]).To(Equal(util.FloatToSimpleString(float64(i) / 100)))
				index += 1
			}
		}
	})

	It("Test redis sorted set command ZRemRangeByRank", func() {
		w.WriteCmdString(handlers.RedisSortedSetCommandZRemRangeByRank, sortedSetCommandTestBaseKey, "1", "2")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("2"))

		w.WriteCmdString(handlers.RedisSortedSetCommandZCard, sortedSetCommandTestBaseKey)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("8"))
	})
})
