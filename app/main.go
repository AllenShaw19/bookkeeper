package app

import (
	"github.com/AllenShaw19/bookkeeper/server/config"
	"github.com/AllenShaw19/bookkeeper/server/proxy"
)

func main() {
	// parse config from filepath
	conf := parseConf()
	// boostrap
	bootstrapDLog(conf)
}

func parseConf() *config.DlogConfig {
	return nil
}

func bootstrapDLog(conf *config.DlogConfig) {
	// start dLogProxy
	dlogProxy := proxy.New(conf)
	dlogProxy.Startup()
}
