package mock

import (
	"fmt"
	"net"
	"testing"
	"redis_go/server"
	"redis_go/protocol"
	"redis_go/tcp"
)

func TestBasicCommand(t *testing.T) {
	// start a mock server
	server := server.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:6379")
	go server.Serve(lis)

	cn, err := net.Dial("tcp", "127.0.0.1:6379")
	defer cn.Close()
	if err != nil {
		t.Fatal(err)
	}

	w := protocol.NewRequestWriter(cn)
	r := protocol.NewResponseReader(cn)

	w.WriteCmdString("PING")
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	responseType, err := r.PeekType()
	if err != nil {
		t.Fatal(err)
	}
	switch responseType {
	case tcp.TypeInline:
		s, _ := r.ReadInlineString()
		fmt.Println(s)
	default:
		t.Fatalf("response type error %+v", responseType)
	}

}
