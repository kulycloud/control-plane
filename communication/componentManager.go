package communication

import (
	"context"
	"errors"
	"fmt"
	protoCommon "github.com/kulycloud/protocol/common"
	"google.golang.org/grpc"
)

var ErrComponentNotFound = errors.New("component not found")

var components = map[string]func(context.Context, *protoCommon.Endpoint) error {
	"route-processor": connectRouteProcessor,
}

var RouteProcessor *RouteProcessorCommunicator = nil

func connectComponent(ctx context.Context, componentType string, endpoint *protoCommon.Endpoint) error {
	fun, ok := components[componentType]
	if !ok {
		return fmt.Errorf("unknown component type %s: %w", componentType, ErrComponentNotFound)
	}
	return fun(ctx, endpoint)
}

func connectRouteProcessor(ctx context.Context, endpoint *protoCommon.Endpoint) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", endpoint.Host, endpoint.Port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not create route-processor connection: %w", err)
	}

	routeProcessor := NewRouteProcessorCommunicator(conn)
	err = routeProcessor.Check(ctx)
	if err != nil {
		return fmt.Errorf("could not ping route-processor: %w", err)
	}

	RouteProcessor = routeProcessor
	logger.Infow("connected route-processor", "endpoint", endpoint)
	return nil
}
