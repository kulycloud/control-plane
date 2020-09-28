package main

import (
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/control-plane/config"
)

func main() {
	logging.Init()
	defer logging.Sync()

	initLogger := logging.GetForComponent("init")

	err := config.ParseConfig()
	if err != nil {
		initLogger.Fatalw("Error parsing config", "error", err)
	}
	initLogger.Infow("Finished parsing config", "config", config.GlobalConfig)
}
