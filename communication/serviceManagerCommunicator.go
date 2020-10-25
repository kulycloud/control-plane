package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"
	protoServices "github.com/kulycloud/protocol/services"
)

var _ commonCommunication.RemoteComponent = &ServiceManagerCommunicator{}

type ServiceManagerCommunicator struct {
	commonCommunication.ComponentCommunicator
	client protoServices.ServiceManagerClient
}

func NewServiceManagerCommunicator(componentCommunicator *commonCommunication.ComponentCommunicator) *ServiceManagerCommunicator {
	return &ServiceManagerCommunicator{ComponentCommunicator: *componentCommunicator, client: protoServices.NewServiceManagerClient(componentCommunicator.GrpcClient)}
}

