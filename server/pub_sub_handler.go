package server

import (
	"regexp"

	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/raw_type"
)

const (
	RedisPubSubCommandPSubscribe   = "PSUBSCRIBE"
	RedisPubSubCommandPublish      = "PUBLISH"
	RedisPubSubCommandPubSub       = "PUBSUB"
	RedisPubSubCommandPUnsubscribe = "PUNSUBSCRIBE"
	RedisPubSubCommandSubscribe    = "SUBSCRIBE"
	RedisPubSubCommandUnsubscribe  = "UNSUBSCRIBE"
)

const (
	PubSubResponseStringMessage      = "message"
	PubSubResponseStringSubscribe    = "subscribe"
	PubSubResponseStringUnsubscribe  = "unsubscribe"
	PubSubResponseStringPSubscribe   = "psubscribe"
	PubSubResponseStringPUnsubscribe = "punsubscribe"
)

type PubSubHandler struct {
}

func (srv *Server) Publish(cli *client.Client) {
	srv.PubSubLock.RLock()
	defer srv.PubSubLock.RUnlock()

	receivers := srv.publishMessage(cli.Argv[1], cli.Argv[2])
	cli.Response(receivers)
}

func (srv *Server) PSubscribe(cli *client.Client) {
	srv.PubSubLock.RLock()
	defer srv.PubSubLock.RUnlock()

	for i := 1; i < len(cli.Argv); i++ {
		srv.subscribePattern(cli, cli.Argv[i])
	}
}

func (srv *Server) PUnsubscribe(cli *client.Client) {
	srv.PubSubLock.RLock()
	defer srv.PubSubLock.RUnlock()

	if len(cli.Argv) == 1 {
		srv.unsubscribeAllPatterns(cli, true)
	} else {
		for i := 1; i < len(cli.Argv); i++ {
			srv.unsubscribePattern(cli, cli.Argv[i], true)
		}
	}
}

// 订阅频道
func (srv *Server) Subscribe(cli *client.Client) {
	srv.PubSubLock.Lock()
	defer srv.PubSubLock.Unlock()

	for i := 1; i < len(cli.Argv); i++ {
		srv.subscribe(cli, cli.Argv[i])
	}
}

// 取消订阅
func (srv *Server) Unsubscribe(cli *client.Client) {
	srv.PubSubLock.Lock()
	defer srv.PubSubLock.Unlock()

	if len(cli.Argv) == 1 {
		srv.unsubscribeAllChannels(cli)
	} else {
		for i := 1; i < len(cli.Argv); i++ {
			srv.unsubscribe(cli, cli.Argv[i])
		}
	}
}

func (srv *Server) PubSub(cli *client.Client) {
	panic("Not implement")
}

// 发送消息
func (srv *Server) publishMessage(channelName string, message string) int {
	var receivers int
	// 发送给普通订阅
	if clients := srv.PubSubChannels[channelName]; clients != nil && clients.ListLength() != 0 {
		iterator := raw_type.ListGetIterator(clients, raw_type.RedisListIteratorDirectionStartHead)
		subClient := iterator.ListNext()
		for subClient != nil {
			responseSlice := make([]interface{}, 3)
			responseSlice[0] = PubSubResponseStringMessage
			responseSlice[1] = channelName
			responseSlice[2] = message
			subClient.NodeValue().(*client.Client).Response(responseSlice)
			subClient.NodeValue().(*client.Client).Flush()

			subClient = iterator.ListNext()
			receivers += 1
		}
	}

	// 发送给模式订阅
	iterator := raw_type.ListGetIterator(srv.PubSubPatterns, raw_type.RedisListIteratorDirectionStartHead)
	item := iterator.ListNext()
	for item != nil {
		pubSubItem := item.NodeValue().(*PubSubPattern)
		if pubSubItem.Exp.Match([]byte(channelName)) {
			responseSlice := make([]interface{}, 4)
			responseSlice[0] = PubSubResponseStringMessage
			responseSlice[1] = pubSubItem.Pattern
			responseSlice[2] = channelName
			responseSlice[3] = message
			pubSubItem.Cli.Response(responseSlice)
			pubSubItem.Cli.Flush()

			receivers += 1
		}
		item = iterator.ListNext()
	}
	return receivers
}

// 为客户端订阅指定频道
func (srv *Server) subscribe(cli *client.Client, channelName string) int {
	var ret int
	if !cli.PubSubChannels.Contains(channelName) {
		ret = 1
		cli.PubSubChannels.Put(channelName, true)
		if _, ok := srv.PubSubChannels[channelName]; !ok {
			srv.PubSubChannels[channelName] = raw_type.ListCreate()
		}
		if node := srv.PubSubChannels[channelName].ListSearchKey(cli); node == nil {
			srv.PubSubChannels[channelName].ListAddNodeTail(cli)
		}
	}
	responseSlice := make([]interface{}, 3)
	responseSlice[0] = PubSubResponseStringSubscribe
	responseSlice[1] = channelName
	responseSlice[2] = cli.PubSubChannels.Size() + cli.PubSubPatterns.ListLength()
	cli.Response(responseSlice)
	cli.Flush()
	return ret
}

// 为客户端取消订阅所有频道
func (srv *Server) unsubscribeAllChannels(cli *client.Client) int {
	var count int
	for channelName := range cli.PubSubChannels.KeySet() {
		count += srv.unsubscribe(cli, channelName.(string))
	}
	cli.PubSubChannels.Clear()
	return count
}

// 为客户端取消订阅指定频道
func (srv *Server) unsubscribe(cli *client.Client, channelName string) int {
	var ret int
	if cli.PubSubChannels.Contains(channelName) {
		ret = 1
		cli.PubSubChannels.RemoveKey(channelName)
		if _, ok := srv.PubSubChannels[channelName]; ok {
			if node := srv.PubSubChannels[channelName].ListSearchKey(cli); node != nil {
				srv.PubSubChannels[channelName].ListRemoveNode(node)
			}
		}
	}
	responseSlice := make([]interface{}, 3)
	responseSlice[0] = PubSubResponseStringUnsubscribe
	responseSlice[1] = channelName
	responseSlice[2] = cli.PubSubChannels.Size() + cli.PubSubPatterns.ListLength()
	cli.Response(responseSlice)
	cli.Flush()
	return ret
}

type PubSubPattern struct {
	Pattern string
	Cli     *client.Client
	Exp     *regexp.Regexp
}

// 为客户端订阅指定模式频道
func (srv *Server) subscribePattern(cli *client.Client, pattern string) int {
	var ret int
	if node := cli.PubSubPatterns.ListSearchKey(pattern); node == nil {
		ret = 1
		cli.PubSubPatterns.ListAddNodeTail(pattern)

		exp, err := regexp.Compile(pattern)
		if err != nil {
			loggers.Errorf("failed to Compile")
			return 0
		}
		srv.PubSubPatterns.ListAddNodeTail(&PubSubPattern{
			Pattern: pattern,
			Cli:     cli,
			Exp:     exp,
		})
	}
	responseSlice := make([]interface{}, 3)
	responseSlice[0] = PubSubResponseStringPSubscribe
	responseSlice[1] = pattern
	responseSlice[2] = cli.PubSubChannels.Size() + cli.PubSubPatterns.ListLength()
	cli.Response(responseSlice)
	cli.Flush()
	return ret
}

// 为客户端取消订阅所有模式频道
func (srv *Server) unsubscribeAllPatterns(cli *client.Client, notifyClient bool) int {
	var count int
	iterator := raw_type.ListGetIterator(srv.PubSubPatterns, raw_type.RedisListIteratorDirectionStartHead)
	item := iterator.ListNext()
	for item != nil {
		pubSubItem := item.NodeValue().(*PubSubPattern)
		nextItem := iterator.ListNext()
		if ret := srv.unsubscribePattern(cli, pubSubItem.Pattern, notifyClient); ret == 1 {
			count += 1
			// 从服务器的pattern中进行移除
			srv.PubSubPatterns.ListRemoveNode(item)
		}
		item = nextItem
	}
	return count
}

// 为客户端取消订阅所有模式频道
func (srv *Server) unsubscribePattern(cli *client.Client, pattern string, notifyClient bool) int {
	var ret int
	if node := cli.PubSubPatterns.ListSearchKey(pattern); node != nil {
		ret = 1
		cli.PubSubPatterns.ListRemoveNode(node)
	}
	if notifyClient {
		responseSlice := make([]interface{}, 3)
		responseSlice[0] = PubSubResponseStringPUnsubscribe
		responseSlice[1] = pattern
		responseSlice[2] = cli.PubSubChannels.Size() + cli.PubSubPatterns.ListLength()
		cli.Response(responseSlice)
		cli.Flush()
	}
	return ret
}
