package proxy

import (
	"github.com/AllenShaw19/bookkeeper/server"
	"github.com/AllenShaw19/bookkeeper/server/config"
)

type DLogProxy struct {
	dLogManager    *DLogManager
	configManager  *ConfigManager
	dLogRpcService server.DLogRpcService
}

func New(conf *config.DlogConfig) *DLogProxy {
	p := &DLogProxy{}
	return p
}

func (p *DLogProxy) Startup() {

}
