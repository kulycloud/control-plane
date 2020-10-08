package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"
	protoRouteProcessor "github.com/kulycloud/protocol/route-processor"
)

var _ commonCommunication.RemoteComponent = &RouteProcessorCommunicator{}
type RouteProcessorCommunicator struct {
	commonCommunication.ComponentCommunicator
	client protoRouteProcessor.RouteProcessorClient
}

func NewRouteProcessorCommunicator(componentCommunicator *commonCommunication.ComponentCommunicator) *RouteProcessorCommunicator {
	return &RouteProcessorCommunicator{ComponentCommunicator: *componentCommunicator, client: protoRouteProcessor.NewRouteProcessorClient(componentCommunicator.GrpcClient)}
}
