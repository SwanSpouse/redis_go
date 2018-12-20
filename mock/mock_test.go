package mock

import (
	"testing"

	"github.com/SwanSpouse/redis_go/conf"
	"github.com/SwanSpouse/redis_go/server"
	"github.com/onsi/ginkgo/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var MockAddr = "127.0.0.1"
var MockPort = 9736

// 如果想要在测试中开启AOF，需要在这里把开关打开
var LoadDataFromAofFile bool

// run all mock test
func TestAll(t *testing.T) {
	p := &server.Program{}
	defaultConfig := conf.NewServerConfig()
	defaultConfig.Port = MockPort
	if LoadDataFromAofFile {
		defaultConfig.AofState = conf.RedisAofOn
	}

	p.InitForMock(defaultConfig)
	p.Start()

	config.GinkgoConfig.FailFast = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Mock Test")

	p.Stop()
}
