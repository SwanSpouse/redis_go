package mock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/server"
)

var _ = Describe("Test Redis key command", func() {
	var cn net.Conn
	var w *RequestWriter
	var r *ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9738")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
		return
	}
	go srv.Serve(lis)
	loggers.Info("redis server start at %s:%s", "127.0.0.1", "9738")

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9738")
		Expect(err).To(BeNil())

		w = NewRequestWriter(cn)
		r = NewResponseReader(cn)
	})

	It("test key command exists and del", func() {
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

		w.WriteCmd(handlers.RedisKeyCommandExists, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmd(handlers.RedisKeyCommandDel, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("1"))

		w.WriteCmd(handlers.RedisStringCommandGet, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("NIL"))

		w.WriteCmd(handlers.RedisKeyCommandExists, []byte(key))
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("0"))
	})

	It("test key command rename", func() {
		key := "my_key"
		value := "my_value"
		w.WriteCmdString(handlers.RedisStringCommandSet, key, value)
		w.Flush()
		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		w.WriteCmdString(handlers.RedisStringCommandGet, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(value))

		newKey := "my_key_new"
		w.WriteCmdString(handlers.RedisKeyCommandRename, key, newKey)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("OK"))

		w.WriteCmdString(handlers.RedisStringCommandGet, key)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal("NIL"))

		w.WriteCmdString(handlers.RedisStringCommandGet, newKey)
		w.Flush()
		ret, err = r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(value))
	})
})
