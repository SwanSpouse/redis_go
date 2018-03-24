package mock

import (
	"fmt"
	"net"
	"redis_go/resp"
	"testing"
)

func TestConnection(t *testing.T) {
	// start a mock server
	server := resp.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:6379")
	go server.Serve(lis)

	cn, err := net.Dial("tcp", "127.0.0.1:6379")
	defer cn.Close()
	if err != nil {
		t.Fatal(err)
	}

	w := resp.NewRequestWriter(cn)
	r := resp.NewResponseReader(cn)

	w.WriteCmdString("PING")
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	responseType, err := r.PeekType()
	if err != nil {
		t.Fatal(err)
	}
	switch responseType {
	case resp.TypeInline:
		s, _ := r.ReadInlineString()
		fmt.Println(s)
	default:
		fmt.Println("No such method")
	}

}
