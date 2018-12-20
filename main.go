package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/SwanSpouse/redis_go/server"
)

// Run runs your Service.
// Run will block until one of the signals specified is received.
func main() {
	prg := &server.Program{}
	prg.Init()
	prg.Start()

	sig := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, sig...)
	<-signalChan

	prg.Stop()
}
