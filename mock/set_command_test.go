package mock

import (
	"fmt"
	"net"

	"github.com/SwanSpouse/redis_go/handlers"
	"github.com/SwanSpouse/redis_go/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestRedisSetCommand", func() {
	var w *RequestWriter
	var r *ResponseReader

	setCmdTestBaseKey := "redis_set_command_test_common_key"

	BeforeEach(func() {
		cn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())

		w = NewRequestWriter(cn)
		r = NewResponseReader(cn)

		// init basic set
		w.WriteCmdString(server.RedisServerCommandFlushAll)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		input := make([]string, 0)
		input = append(input, setCmdTestBaseKey)
		for i := 0; i < 10; i++ {
			input = append(input, fmt.Sprintf("item%d", i))
		}
		w.WriteCmdString(handlers.RedisSetCommandSADD, input...)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))
	})

	It("Test redis set command SAdd SCard", func() {
		w.WriteCmdString(handlers.RedisSetCommandSADD, setCmdTestBaseKey, "value")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("11"))

		w.WriteCmdString(handlers.RedisSetCommandSADD, setCmdTestBaseKey, "item1")
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("0"))

		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("11"))
	})

	It("Test redis set command SIsMember SMembers", func() {
		for i := 0; i < 10; i += 1 {
			w.WriteCmdString(handlers.RedisSetCommandSISMEMBER, setCmdTestBaseKey, fmt.Sprintf("item%d", i))
			w.Flush()
			ret, _ := r.Read()
			Expect(ret[0]).To(Equal("1"))

			w.WriteCmdString(handlers.RedisSetCommandSISMEMBER, setCmdTestBaseKey, fmt.Sprintf("value%d", i))
			w.Flush()
			ret, _ = r.Read()
			Expect(ret[0]).To(Equal("0"))
		}

		members := make(map[string]bool)
		for i := 0; i < 10; i += 1 {
			members[fmt.Sprintf("item%d", i)] = true
		}
		w.WriteCmdString(handlers.RedisSetCommandSMEMBERS, setCmdTestBaseKey)
		w.Flush()
		ret, _ := r.Read()
		Expect(len(ret)).To(Equal(len(members)))
		for _, item := range ret {
			if _, ok := members[item]; !ok {
				Expect(true).To(Equal(false))
			}
		}
	})

	It("Test redis set command SPop SRem SRandMember", func() {
		members := make(map[string]bool)
		for i := 0; i < 10; i += 1 {
			members[fmt.Sprintf("item%d", i)] = true
		}
		w.WriteCmdString(handlers.RedisSetCommandSPOP, setCmdTestBaseKey)
		w.Flush()
		ret, _ := r.Read()
		if _, ok := members[ret[0]]; !ok {
			Expect(true).To(Equal(false))
		}
		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("9"))

		w.WriteCmdString(handlers.RedisSetCommandSRANDMEMBER, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		if _, ok := members[ret[0]]; !ok {
			Expect(true).To(Equal(false))
		}
		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("9"))

		w.WriteCmdString(handlers.RedisSetCommandSREM, setCmdTestBaseKey, "item1")
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("1"))
		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("8"))

		w.WriteCmdString(handlers.RedisSetCommandSREM, setCmdTestBaseKey, "value1")
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("0"))
		w.WriteCmdString(handlers.RedisSetCommandSCARD, setCmdTestBaseKey)
		w.Flush()
		ret, _ = r.Read()
		Expect(ret[0]).To(Equal("8"))
	})
})
