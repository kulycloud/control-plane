package communication

import (
	"context"
	protoCommon "github.com/kulycloud/protocol/common"
	protoRouteProcessor "github.com/kulycloud/protocol/route-processor"
	"google.golang.org/grpc"
)

type RouteProcessorCommunicator struct {
	client protoRouteProcessor.RouteProcessorClient
}

func NewRouteProcessorCommunicator(grpcClient grpc.ClientConnInterface) *RouteProcessorCommunicator {
	return &RouteProcessorCommunicator{client: protoRouteProcessor.NewRouteProcessorClient(grpcClient)}
}

func (communicator *RouteProcessorCommunicator) Check(ctx context.Context) error {
	_, err := communicator.client.Ping(ctx, &protoCommon.Empty{})
	return err
}
