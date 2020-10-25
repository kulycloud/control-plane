package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"
	protoIngress "github.com/kulycloud/protocol/ingress"
)

var _ commonCommunication.RemoteComponent = &IngressCommunicator{}

type IngressCommunicator struct {
	commonCommunication.ComponentCommunicator
	client protoIngress.IngressClient
}

func NewIngressCommunicator(componentCommunicator *commonCommunication.ComponentCommunicator) *IngressCommunicator {
	return &IngressCommunicator{ComponentCommunicator: *componentCommunicator, client: protoIngress.NewIngressClient(componentCommunicator.GrpcClient)}
}
