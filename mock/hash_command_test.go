package mock_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/protocol"
	"redis_go/server"
)

var _ = FDescribe("TestRedisHashCommand", func() {
	var cn net.Conn
	var w *protocol.RequestWriter
	var r *protocol.ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9733")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
	}
	go srv.Serve(lis)
	loggers.Info("redis server start at %s:%s", "127.0.0.1", "9733")


	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9733")
		w = protocol.NewRequestWriter(cn)
		r = protocol.NewResponseReader(cn)
	})

	It("Test redis hash command HSet, HGet HSetNx HExists HDel HLen", func() {
		key := "redis_hash_command_test_common_key"
		w.WriteCmdString(handlers.RedisHashCommandHSet, key, "key", "value")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHGet, key, "key")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("value"))

		w.WriteCmdString(handlers.RedisHashCommandHLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHExists, key, "key")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHDel, key, "key")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))

		w.WriteCmdString(handlers.RedisHashCommandHExists, key, "key")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))

		w.WriteCmdString(handlers.RedisHashCommandHSetNX, key, "key", "value")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHSetNX, key, "key", "value")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))
	})

	It("Test redis hash command HMSet HMGet HKeys HVals", func() {
		key := "redis_hash_command_test_common_key"
		input := make([]string, 0)
		input = append(input, key)
		for i := 0; i < 10; i++ {
			input = append(input, fmt.Sprintf("key%d", i))
			input = append(input, fmt.Sprintf("value%d", i))
		}
		w.WriteCmdString(handlers.RedisHashCommandHMSet, input...)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisHashCommandHGet, key, fmt.Sprintf("key%d", i))
			w.Flush()
			ret, err = r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("value%d", i)))
		}
		w.WriteCmdString(handlers.RedisHashCommandHGetAll, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())

		retMapHGetAll := make(map[string]string)
		for i := 0; i < len(ret); i += 2 {
			retMapHGetAll[ret[i]] = ret[i+1]
		}
		for i := 0; i < 10; i += 1 {
			if value, ok := retMapHGetAll[fmt.Sprintf("key%d", i)]; ok {
				Expect(value).To(Equal(fmt.Sprintf("value%d", i)))
			} else {
				// 不应该到达的分支
				Expect(true).To(Equal(false))
			}
		}

		retMapHKeys := make(map[string]bool)
		w.WriteCmdString(handlers.RedisHashCommandHKeys, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		for _, item := range ret {
			retMapHKeys[item] = true
		}
		for i := 0; i < 10; i++ {
			if value, ok := retMapHKeys[fmt.Sprintf("key%d", i)]; ok {
				Expect(value).To(Equal(true))
			} else {
				// 不应该到达的分支
				Expect(true).To(Equal(false))
			}
		}

		retMapHVals := make(map[string]bool)
		w.WriteCmdString(handlers.RedisHashCommandHVals, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		for _, item := range ret {
			retMapHVals[item] = true
		}

		for i := 0; i < 10; i++ {
			if value, ok := retMapHVals[fmt.Sprintf("value%d", i)]; ok {
				Expect(value).To(Equal(true))
			} else {
				// 不应该到达的分支
				Expect(true).To(Equal(false))
			}
		}
	})

	It("Test redis hash command HIncrBy, HIncrByFloat", func() {
		key := "redis_hash_command_test_common_key_2"
		w.WriteCmdString(handlers.RedisHashCommandHIncrBy, key, "one", "1")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHIncrBy, key, "one", "10")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("11"))

		w.WriteCmdString(handlers.RedisHashCommandHIncrBy, key, "one", "-10")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisHashCommandHIncrByFloat, key, "float", "-10")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("-10.00"))

		w.WriteCmdString(handlers.RedisHashCommandHIncrByFloat, key, "float", "10.999")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0.999"))
	})

})
