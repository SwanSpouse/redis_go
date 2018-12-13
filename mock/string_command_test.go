package mock

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/server"
)

var _ = Describe("TestRedisStringCommand", func() {
	var cn net.Conn
	var w *RequestWriter
	var r *ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9731")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
	}
	go srv.Serve(lis)
	loggers.Info("redis server start at %s:%s", "127.0.0.1", "9731")

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9731")
		w = NewRequestWriter(cn)
		r = NewResponseReader(cn)
	})

	It("test redis string command set and get", func() {
		key := "my_key"
		value := "my_value"
		w.WriteCmd(handlers.RedisStringCommandSet, []byte(key), []byte(value))
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		w.WriteCmd(handlers.RedisStringCommandGet, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(value))
	})

	It("test redis string command mset and mget", func() {
		input := make([][]byte, 0)
		for i := 0; i < 10; i++ {
			input = append(input, []byte(fmt.Sprintf("key%d", i)))
			input = append(input, []byte(fmt.Sprintf("value%d", i)))
		}
		w.WriteCmd(handlers.RedisStringCommandMSet, input...)
		w.Flush()

		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		input = make([][]byte, 0)
		for i := 0; i < 10; i++ {
			input = append(input, []byte(fmt.Sprintf("key%d", i)))
		}
		w.WriteCmd(handlers.RedisStringCommandMGet, input...)
		w.Flush()

		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(len(ret)).To(Equal(10))
		for i := 0; i < 10; i++ {
			Expect(ret[i]).To(Equal(fmt.Sprintf("value%d", i)))
		}
	})

	// TODO RDB暂时没有考虑编码的问题。所以暂时关闭这个测试
	//It("test redis string type and encodings", func() {
	//	key := "number"
	//	value := "123"
	//	w.WriteCmd(handlers.RedisStringCommandSet, []byte(key), []byte(value))
	//	w.Flush()
	//	ret, err := r.Read()
	//	Expect(err).To(BeNil())
	//	Expect(ret[0]).To(Equal("OK"))
	//
	//	w.WriteCmd(handlers.RedisKeyCommandType, []byte(key))
	//	w.Flush()
	//	ret, err = r.Read()
	//	Expect(err).To(BeNil())
	//	Expect(ret[0]).To(Equal(encodings.RedisTypeString))
	//
	//	w.WriteCmd(handlers.RedisKeyCommandObject, []byte("encoding"), []byte(key))
	//	w.Flush()
	//	ret, err = r.Read()
	//	Expect(err).To(BeNil())
	//	Expect(ret[0]).To(Equal(encodings.RedisEncodingInt))
	//
	//	w.WriteCmd(handlers.RedisStringCommandAppend, []byte(key), []byte(" "))
	//	w.Flush()
	//	ret, err = r.Read()
	//	Expect(err).To(BeNil())
	//	Expect(ret[0]).To(Equal(fmt.Sprintf("%d", len(value+"x"))))
	//
	//	w.WriteCmd(handlers.RedisKeyCommandObject, []byte("encoding"), []byte(key))
	//	w.Flush()
	//	ret, err = r.Read()
	//	Expect(err).To(BeNil())
	//	Expect(ret[0]).To(Equal(encodings.RedisEncodingRaw))
	//})

	It("test redis string incr and decr", func() {
		key := "number"
		value := "123"
		w.WriteCmd(handlers.RedisStringCommandSet, []byte(key), []byte(value))
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		w.WriteCmd(handlers.RedisStringCommandIncr, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("124"))

		w.WriteCmd(handlers.RedisStringCommandDecr, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("123"))

		w.WriteCmd(handlers.RedisStringCommandIncrBy, []byte(key), []byte("10"))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("133"))

		w.WriteCmd(handlers.RedisStringCommandDecrBy, []byte(key), []byte("10"))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("123"))

		w.WriteCmd(handlers.RedisStringCommandIncrByFloat, []byte(key), []byte("10.12"))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("133.12"))

		w.WriteCmd(handlers.RedisStringCommandStrLen, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("6"))
	})
})
