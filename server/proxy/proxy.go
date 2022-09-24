package proxy

import (
	"github.com/AllenShaw19/bookkeeper/server"
)

type DLogProxy struct {
	dLogManager    *DLogManager
	configManager  *ConfigManager
	dLogRpcService server.DLogRpcService
}

func New(conf *server.DLogConfig) *DLogProxy {
	p := &DLogProxy{}
	return p
}

func (p *DLogProxy) Startup() {

}
