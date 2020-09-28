module github.com/kulycloud/control-plane

go 1.15

require (
	github.com/kulycloud/common v1.0.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
)

replace github.com/kulycloud/common v1.0.0 => ../common
