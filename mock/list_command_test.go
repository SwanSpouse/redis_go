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

var _ = FDescribe("TestRedisListCommand", func() {
	var cn net.Conn
	var w *protocol.RequestWriter
	var r *protocol.ResponseReader

	commonKey := "redis_list_test_common_key"
	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9732")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
	}
	go srv.Serve(lis)
	loggers.Info("redis server start at %s:%s", "127.0.0.1", "9732")

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9732")
		w = protocol.NewRequestWriter(cn)
		r = protocol.NewResponseReader(cn)

		input := make([]string, 0)
		input = append(input, commonKey)
		for i := 0; i < 10; i++ {
			input = append(input, fmt.Sprintf("value%d", i))
		}
		w.WriteCmdString(handlers.RedisListCommandRPush, input...)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))
	})

	AfterEach(func() {
		w.WriteCmdString(handlers.RedisKeyCommandDel, commonKey)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))
	})

	It("test redis list command LPush LPop LLen", func() {
		key := "test_redis_list_command_lpush"
		w.WriteCmdString(handlers.RedisListCommandLPush, key, "value")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisListCommandLLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisListCommandLPush, key, "value2")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("2"))

		w.WriteCmdString(handlers.RedisListCommandLLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("2"))

		w.WriteCmdString(handlers.RedisListCommandLPop, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("value2"))

		w.WriteCmdString(handlers.RedisListCommandLLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisListCommandLPop, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("value"))

		w.WriteCmdString(handlers.RedisListCommandLLen, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))
	})

	It("test redis list command RPush RPop", func() {
		key := "test_redis_list_command_rpush"
		w.WriteCmdString(handlers.RedisListCommandRPush, key, "value")
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisListCommandRPush, key, "value2")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("2"))

		w.WriteCmdString(handlers.RedisListCommandRPop, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("value2"))

		w.WriteCmdString(handlers.RedisListCommandRPop, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("value"))
	})

	It("test redis list command LIndex and LInsert ", func() {
		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisListCommandLIndex, commonKey, fmt.Sprintf("%d", i))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("value%d", i)))
		}

		for i := -1; i >= -10; i-- {
			w.WriteCmdString(handlers.RedisListCommandLIndex, commonKey, fmt.Sprintf("%d", i))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("value%d", 10+i)))
		}

		valueNewInsert := "valueNewInsert"
		w.WriteCmdString(handlers.RedisListCommandLInsert, commonKey, "BEFORE", "value0", valueNewInsert)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("11"))
	})

	It("test redis list command LSet ", func() {
		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisListCommandLSet, commonKey, fmt.Sprintf("%d", i), fmt.Sprintf("newValue%d", i))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("OK"))
		}

		for i := 0; i < 10; i++ {
			w.WriteCmdString(handlers.RedisListCommandLIndex, commonKey, fmt.Sprintf("%d", i))
			w.Flush()
			ret, err := r.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(fmt.Sprintf("newValue%d", i)))
		}
	})

	It("test redis list command LRem ", func() {
		curKey := "redis_list_test_remove_key"
		input := make([]string, 0)
		input = append(input, curKey)
		for i := 0; i < 10; i++ {
			input = append(input, fmt.Sprintf("value%d", i/5))
		}
		w.WriteCmdString(handlers.RedisListCommandLPush, input...)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("10"))

		w.WriteCmdString(handlers.RedisListCommandLRem, curKey, "4", "value0")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("4"))

		w.WriteCmdString(handlers.RedisListCommandLLen, curKey)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("6"))

		w.WriteCmdString(handlers.RedisListCommandLRem, curKey, "4", "value0")
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmdString(handlers.RedisListCommandLLen, curKey)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("5"))

		w.WriteCmdString(handlers.RedisKeyCommandDel, curKey)
		w.Flush()
	})
})
