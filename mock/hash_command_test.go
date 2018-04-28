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

func TestRedisHashCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Redis Hash Commands")
}

var _ = Describe("TestRedisHashCommand", func() {
	var cn net.Conn
	var w *protocol.RequestWriter
	var r *protocol.ResponseReader

	srv := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:9737")
	if err != nil {
		loggers.Errorf("server start error %+v", err)
	}
	go srv.Serve(lis)

	BeforeEach(func() {
		cn, err = net.Dial("tcp", "127.0.0.1:9737")
		w = protocol.NewRequestWriter(cn)
		r = protocol.NewResponseReader(cn)
	})

	AfterEach(func() {
		cn.Close()
	})

	It("Test redis hash command HSet, HGet", func() {
		key := "redis_hash_command_test_common_key"
		w.WriteCmdString(handlers.RedisHashCommandHSet, key, "key", "value")
		w.Flush()

		ret, err := r.Read()
		Ω(err).To(BeNil())
		Ω(ret[0]).To(Equal("1"))
	})
})
