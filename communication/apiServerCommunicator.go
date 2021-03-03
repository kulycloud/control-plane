package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"
	protoApiServer "github.com/kulycloud/protocol/api-server"
)

var _ commonCommunication.RemoteComponent = &ApiServerCommunicator{}

type ApiServerCommunicator struct {
	commonCommunication.ComponentCommunicator
	client protoApiServer.ApiServerClient
}

func NewApiServerCommunicator(componentCommunicator *commonCommunication.ComponentCommunicator) *ApiServerCommunicator {
	return &ApiServerCommunicator{ComponentCommunicator: *componentCommunicator, client: protoApiServer.NewApiServerClient(componentCommunicator.GrpcClient)}
}
