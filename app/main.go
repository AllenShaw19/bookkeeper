package app

import (
	"github.com/AllenShaw19/bookkeeper/server"
	"github.com/AllenShaw19/bookkeeper/server/proxy"
)

func main() {
	// parse config from filepath
	conf := parseConf()
	// boostrap
	bootstrapDLog(conf)
}

func parseConf() *server.DLogConfig {
	return nil
}

func bootstrapDLog(conf *server.DLogConfig) {
	// start dLogProxy
	dlogProxy := proxy.New(conf)
	dlogProxy.Startup()
}
