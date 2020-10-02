package communication

import (
	"context"
	"fmt"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/control-plane/config"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
	"google.golang.org/grpc"
	"net"
)


var _ protoControlPlane.ControlPlaneServer = &Listener{}

var logger = logging.GetForComponent("communication")

type Listener struct {
	server *grpc.Server
	listener net.Listener
}

func NewListener() *Listener {
	return &Listener{}
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

func (listener *Listener) RegisterComponent(ctx context.Context, message *protoControlPlane.RegisterComponentRequest) (*protoControlPlane.RegisterComponentResult, error) {
	logger.Infow("registering component", "component", message.Type)
	err := connectComponent(ctx, message.Type, message.Endpoint)
	if err != nil {
		logger.Warnw("error connecting to component", "type", message.Type, "endpoint", message.Endpoint, "error", err)
	}

	return &protoControlPlane.RegisterComponentResult{}, err
}
