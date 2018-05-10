package rdb

import (
	"fmt"
	"github.com/cupcake/rdb/nopdecoder"
	"gopkg.in/bufio.v1"
	"os"
	"testing"
)

type decoder struct {
	db int
	i  int
	nopdecoder.NopDecoder
}

func (p *decoder) StartDatabase(n int) {
	p.db = n
}

func (p *decoder) Set(key, value []byte, expiry int64) {
	fmt.Printf("db=%d %q -> %q\n", p.db, key, value)
}

func (p *decoder) Hset(key, field, value []byte) {
	fmt.Printf("db=%d %q . %q -> %q\n", p.db, key, field, value)
}

func (p *decoder) Sadd(key, member []byte) {
	fmt.Printf("db=%d %q { %q }\n", p.db, key, member)
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *decoder) Rpush(key, value []byte) {
	fmt.Printf("db=%d %q[%d] -> %q\n", p.db, key, p.i, value)
	p.i++
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *decoder) Zadd(key []byte, score float64, member []byte) {
	fmt.Printf("db=%d %q[%d] -> {%q, score=%g}\n", p.db, key, p.i, member, score)
	p.i++
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func TestReadFile(t *testing.T) {
	f, err := os.Open("./dump.string_obj.rdb")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	count := 1
	for {
		flag, err := reader.ReadByte()
		if err != nil {
			break
		}
		fmt.Printf("%dth byte ==>%+v %U %s\n", count, flag, flag, string(flag))
		count += 1
	}

}