package mock

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/protocol"
	"redis_go/server"
	"redis_go/tcp"
	"testing"
)

func TestRedisClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Redis client")
}

var _ = Describe("MockRedisClient", func() {
	var cn net.Conn
	var w *protocol.RequestWriter
	var r *protocol.ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9739")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
		return
	}

	go srv.Serve(lis)

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9739")
		Expect(err).To(BeNil())

		w = protocol.NewRequestWriter(cn)
		r = protocol.NewResponseReader(cn)
	})

	AfterEach(func() {
		cn.Close()
	})

	It("test ping command", func() {
		w.WriteCmdString("PING")
		err := w.Flush()
		Expect(err).To(BeNil())

		responseType, err := r.PeekType()
		Expect(err).To(BeNil())

		switch responseType {
		case tcp.TypeInline:
			s, _ := r.ReadInlineString()
			Expect(s).To(Equal("PONG"))
		case tcp.TypeBulk:
			s, _ := r.ReadBulkString()
			Expect(s).To(Equal("PONG"))
		default:
			panic(fmt.Sprintf("response type error %+v", responseType))
		}
	})

	It("test unknown command", func() {
		w.WriteRawString("*1\r\n$3\r\nlol\r\n")
		err := w.Flush()
		Expect(err).To(BeNil())
		Expect(r.PeekType()).To(Equal(tcp.TypeError))

		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(ContainSubstring("unknown command"))
	})

	It("test input err invalid bulk length", func() {
		w.WriteRawString("*3\r\n$3\r\nget\r\n$-1\r\n$3\r\nbar\r\n")
		err := w.Flush()
		Expect(err).To(BeNil())

		Expect(r.PeekType()).To(Equal(tcp.TypeError))

		ret, err := r.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(re.ErrInvalidBulkLength.Error()))
	})
})
