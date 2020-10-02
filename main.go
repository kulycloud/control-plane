package main

import (
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/control-plane/communication"
	"github.com/kulycloud/control-plane/config"
)

func main() {
	initLogger := logging.GetForComponent("init")
	defer logging.Sync()

	err := config.ParseConfig()
	if err != nil {
		initLogger.Fatalw("Error parsing config", "error", err)
	}
	initLogger.Infow("Finished parsing config", "config", config.GlobalConfig)

	initLogger.Info("Starting listener")
	listener := communication.NewListener()
	err = listener.Start()
	if err != nil {
		initLogger.Panicw("error initializing listener", "error", err)
	}
}
