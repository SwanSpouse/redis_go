package rdb

import (
	"os"
	"testing"
	"fmt"
)

func TestWriteFile(t *testing.T) {
	f, err := os.Create("./output3")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	redis := []byte("REDIS0006")

	fmt.Printf("%+v\n", redis)

	f.Write(redis)
	f.Sync()
}
