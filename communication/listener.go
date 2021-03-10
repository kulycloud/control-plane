package communication

import (
	"context"
	"fmt"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/control-plane/config"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
    protoCommon "github.com/kulycloud/protocol/common"
	"google.golang.org/grpc"
	"net"
)

var _ protoControlPlane.ControlPlaneServer = &Listener{}

var logger = logging.GetForComponent("communication")

type Listener struct {
	protoControlPlane.UnimplementedControlPlaneServer
	server   *grpc.Server
	listener net.Listener

    eventStreams map[string]*eventStream
}

func NewListener() *Listener {
	return &Listener{
        eventStreams: make(map[string]*eventStream, 0),
    }
}

func (listener *Listener) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.GlobalConfig.RPCPort))
	if err != nil {
		return err
	}
	listener.listener = lis
	listener.server = grpc.NewServer()
	protoControlPlane.RegisterControlPlaneServer(listener.server, listener)
	logger.Infow("serving", "port", config.GlobalConfig.RPCPort)
	return listener.server.Serve(listener.listener)
}

func (listener *Listener) RegisterComponent(message *protoControlPlane.RegisterComponentRequest, stream protoControlPlane.ControlPlane_RegisterComponentServer) error {
	logger.Infow("registering component", "type", message.Type)
	err := GlobalComponentManager.ConnectComponent(stream.Context(), message.Type, message.Endpoint)
	if err != nil {
		logger.Warnw("error connecting to component", "type", message.Type, "endpoint", message.Endpoint, "error", err)
	}
    eventStream := newEventStream(message.Endpoint, stream)
    listener.eventStreams[eventStream.destination] = eventStream

	return err
}

func (listener *Listener) CreateEvent(ctx context.Context, event *protoCommon.Event) (*protoCommon.Empty, error) {
    for destination, stream := range listener.eventStreams {
        err := stream.send(event)
        if err != nil {
            logger.Warn("error while sending event to stream %w", err)
            delete(listener.eventStreams, destination)
        }
    }
    return nil, nil
}

func (listener *Listener) ListenToEvent(ctx context.Context, request *protoControlPlane.ListenToEventRequest) (*protoCommon.Empty, error) {
    
    return nil, nil
}

