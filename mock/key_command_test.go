package mock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/protocol"
	"redis_go/server"
	"testing"
)

func TestRedisKeyCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Redis key command")
}

var _ = Describe("Test Redis key command", func() {
	var cn net.Conn
	var w *protocol.RequestWriter
	var r *protocol.ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9730")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
		return
	}

	go srv.Serve(lis)

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9730")
		Expect(err).To(BeNil())

		w = protocol.NewRequestWriter(cn)
		r = protocol.NewResponseReader(cn)
	})

	AfterEach(func() {
		cn.Close()
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
})