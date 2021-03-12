package communication

import (
	"context"
	"fmt"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/control-plane/config"
	protoCommon "github.com/kulycloud/protocol/common"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var _ protoControlPlane.ControlPlaneServer = &Listener{}

var logger = logging.GetForComponent("communication")

type Listener struct {
	protoControlPlane.UnimplementedControlPlaneServer
	server   *grpc.Server
	listener net.Listener

	streamMutex  sync.Mutex
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

	// create event stream and store it
	eventStream := newEventStream(message.Endpoint, stream)
	logger.Debugf("created new event stream for %s", eventStream.destination)
	listener.streamMutex.Lock()
	listener.eventStreams[eventStream.destination] = eventStream
	listener.streamMutex.Unlock()

	// send confirmation that the stream has been created
	// no errors indicate a working stream so it waits until the streams context is closed
	err = eventStream.sendConfirmation()
	if err == nil {
		<-stream.Context().Done()
		logger.Debug("stream context done, deleting stream")
		listener.streamMutex.Lock()
		delete(listener.eventStreams, eventStream.destination)
		listener.streamMutex.Unlock()
	}

	return err
}

func (listener *Listener) CreateEvent(ctx context.Context, event *protoCommon.Event) (*protoCommon.Empty, error) {
	logger.Debugf("creating event %v", event)
	listener.streamMutex.Lock()
	defer listener.streamMutex.Unlock()
	for destination, stream := range listener.eventStreams {
		err := stream.send(event)
		if err != nil {
			logger.Warnf("error while sending event to stream %v, deleting stream", err)
			delete(listener.eventStreams, destination)
		}
		logger.Debugf("sent event to %s", destination)
	}
	return &protoCommon.Empty{}, nil
}

func (listener *Listener) ListenToEvent(ctx context.Context, request *protoControlPlane.ListenToEventRequest) (*protoCommon.Empty, error) {
	logger.Debugf("%s requested to receive %s events", request.Destination, request.Type)
	listener.streamMutex.Lock()
	defer listener.streamMutex.Unlock()
	stream, ok := listener.eventStreams[request.Destination]
	if !ok {
		return &protoCommon.Empty{}, fmt.Errorf("could not find associated stream")
	}
	stream.listenOnEvent(request.Type)
	return &protoCommon.Empty{}, nil
}
