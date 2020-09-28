package config

import (
	commonConfig "github.com/kulycloud/common/config"
)

type Config struct {
	RPCPort uint16 `configName:"rpcPort" defaultValue:"12270"`
}

var GlobalConfig = &Config{}

func ParseConfig() error {
	parser := commonConfig.NewParser()
	parser.AddProvider(commonConfig.NewCliParamProvider())
	parser.AddProvider(commonConfig.NewEnvironmentVariableProvider())

	return parser.Populate(GlobalConfig)
}
