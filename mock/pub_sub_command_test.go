package mock

import (
	"fmt"
	"net"

	"github.com/SwanSpouse/redis_go/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestPubSubCommand", func() {
	var w *RequestWriter
	var r *ResponseReader

	var normalChannelName = "redis_pub_sub_command_key_normal"
	var patternChannelName = "redis_pub_sub_command_key_pattern.*"

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
	})

	It("test redis pub/sub command", func() {
		// 创建两个新的客户端来接收消息
		SubscriberCn1, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())
		Subscriber1Writer := NewRequestWriter(SubscriberCn1)
		Subscriber1Reader := NewResponseReader(SubscriberCn1)

		// 首先订阅normalChannelName频道，并验证返回值
		Subscriber1Writer.WriteCmdString(server.RedisPubSubCommandSubscribe, normalChannelName)
		Subscriber1Writer.Flush()

		ret, err := Subscriber1Reader.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(server.PubSubResponseStringSubscribe))
		Expect(ret[1]).To(Equal(normalChannelName))
		Expect(ret[2]).To(Equal("1"))

		SubscriberCn2, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())
		Subscriber2Writer := NewRequestWriter(SubscriberCn2)
		Subscriber2Reader := NewResponseReader(SubscriberCn2)

		// 首先订阅normalChannelName频道，并验证返回值
		Subscriber2Writer.WriteCmdString(server.RedisPubSubCommandSubscribe, normalChannelName)
		Subscriber2Writer.Flush()

		ret, err = Subscriber2Reader.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(server.PubSubResponseStringSubscribe))
		Expect(ret[1]).To(Equal(normalChannelName))
		Expect(ret[2]).To(Equal("1"))

		// 创建一个新的客户端来发送消息
		PublishCn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())

		PublisherWriter := NewRequestWriter(PublishCn)
		PublisherReader := NewResponseReader(PublishCn)

		messageContent := "This is message :%d"
		for i := 0; i < 10; i++ {
			// 发送消息，并验证收到的结果
			PublisherWriter.WriteCmdString(server.RedisPubSubCommandPublish, normalChannelName, fmt.Sprintf(messageContent, i))
			PublisherWriter.Flush()

			ret, err := PublisherReader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("2"))

			// 验证订阅客户端收到的消息
			ret, err = Subscriber1Reader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(server.PubSubResponseStringMessage))
			Expect(ret[1]).To(Equal(normalChannelName))
			Expect(ret[2]).To(Equal(fmt.Sprintf(messageContent, i)))

			ret, err = Subscriber2Reader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal(server.PubSubResponseStringMessage))
			Expect(ret[1]).To(Equal(normalChannelName))
			Expect(ret[2]).To(Equal(fmt.Sprintf(messageContent, i)))
		}

		SubscriberCn1.Close()
		SubscriberCn2.Close()
		PublishCn.Close()
	})

	It("test redis ppub/psub command", func() {
		// 创建两个新的客户端来接收消息
		SubscriberCn1, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())
		Subscriber1Writer := NewRequestWriter(SubscriberCn1)
		Subscriber1Reader := NewResponseReader(SubscriberCn1)

		// 首先订阅patternChannelName频道，并验证返回值
		Subscriber1Writer.WriteCmdString(server.RedisPubSubCommandPSubscribe, patternChannelName)
		Subscriber1Writer.Flush()

		ret, err := Subscriber1Reader.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(server.PubSubResponseStringPSubscribe))
		Expect(ret[1]).To(Equal(patternChannelName))
		Expect(ret[2]).To(Equal("1"))

		SubscriberCn2, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())
		Subscriber2Writer := NewRequestWriter(SubscriberCn2)
		Subscriber2Reader := NewResponseReader(SubscriberCn2)

		// 首先订阅patternChannelName频道，并验证返回值
		Subscriber2Writer.WriteCmdString(server.RedisPubSubCommandPSubscribe, patternChannelName)
		Subscriber2Writer.Flush()

		ret, err = Subscriber2Reader.Read()
		Expect(err).To(BeNil())
		Expect(ret[0]).To(Equal(server.PubSubResponseStringPSubscribe))
		Expect(ret[1]).To(Equal(patternChannelName))
		Expect(ret[2]).To(Equal("1"))

		// 创建一个新的客户端来发送消息
		PublishCn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", MockAddr, MockPort))
		Expect(err).To(BeNil())

		PublisherWriter := NewRequestWriter(PublishCn)
		PublisherReader := NewResponseReader(PublishCn)

		messageContent := "This is message :%d"
		for i := 0; i < 10; i++ {
			targetChannelName := fmt.Sprintf("%sasdlfkj%d", patternChannelName, i)
			// 发送消息，并验证收到的结果
			PublisherWriter.WriteCmdString(server.RedisPubSubCommandPublish, targetChannelName, fmt.Sprintf(messageContent, i))
			PublisherWriter.Flush()

			ret, err := PublisherReader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("2"))

			// 验证订阅客户端收到的消息
			ret, err = Subscriber1Reader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("message"))
			Expect(ret[1]).To(Equal(patternChannelName))
			Expect(ret[2]).To(Equal(targetChannelName))
			Expect(ret[3]).To(Equal(fmt.Sprintf(messageContent, i)))

			ret, err = Subscriber2Reader.Read()
			Expect(err).To(BeNil())
			Expect(ret[0]).To(Equal("message"))
			Expect(ret[1]).To(Equal(patternChannelName))
			Expect(ret[2]).To(Equal(targetChannelName))
			Expect(ret[3]).To(Equal(fmt.Sprintf(messageContent, i)))
		}

		SubscriberCn1.Close()
		SubscriberCn2.Close()
		PublishCn.Close()
	})
})
