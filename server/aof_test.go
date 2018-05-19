package server

import (
	"redis_go/loggers"
	"testing"
)

func TestCatAppendOnlyGenericCommand(t *testing.T) {
	buf := make([]byte, 0)
	argc := 3
	argv := []string{"SET", "LMJ", "123"}
	ret := catAppendOnlyGenericCommand(buf, argc, argv)
	loggers.Info("ret:%s", ret)
}
