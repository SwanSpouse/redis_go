package mock

import (
	"fmt"
	"net"
	"redis_go/resp"
	"testing"
)

func TestConnection(t *testing.T) {
	// start a mock server
	server := networking.NewServer(nil)
	lis, err := net.Listen("tcp", "127.0.0.1:6379")
	go server.Serve(lis)

	cn, err := net.Dial("tcp", "127.0.0.1:6379")
	defer cn.Close()
	if err != nil {
		t.Fatal(err)
	}

	w := networking.NewRequestWriter(cn)
	r := networking.NewResponseReader(cn)

	w.WriteCmdString("PING")
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	responseType, err := r.PeekType()
	if err != nil {
		t.Fatal(err)
	}
	switch responseType {
	case networking.TypeInline:
		s, _ := r.ReadInlineString()
		fmt.Println(s)
	default:
		fmt.Println("No such method")
	}

}
