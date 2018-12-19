package mock

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"redis_go/conf"
	"redis_go/server"
	"testing"
)

var MockAddr = "127.0.0.1"
var MockPort = 9736

// run all mock test
func TestAll(t *testing.T) {
	p := &server.Program{}
	defaultConfig := conf.NewServerConfig()
	defaultConfig.Port = MockPort
	defaultConfig.AofState = conf.RedisAofOff

	p.InitForMock(defaultConfig)
	p.Start()

	config.GinkgoConfig.FailFast = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Mock Test")

	p.Stop()
}
