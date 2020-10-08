package communication

import (
	"context"
	"errors"
	"fmt"
	commonCommunication "github.com/kulycloud/common/communication"
	commonStorage "github.com/kulycloud/common/storage"
	protoCommon "github.com/kulycloud/protocol/common"
	"google.golang.org/grpc"
)

var ErrComponentNotFound = errors.New("component not found")

type RegisterHandler = func(context.Context, *ComponentManager, string, commonCommunication.RemoteComponent, *protoCommon.Endpoint)
type ComponentFactory = func(context.Context, *ComponentManager, *commonCommunication.ComponentCommunicator) (commonCommunication.RemoteComponent, error)

type ComponentManager struct {
	GeneralRegisterHandlers	[]RegisterHandler
	RegisterHandlers map[string][]RegisterHandler
	Components []commonCommunication.RemoteComponent
	factorySetters map[string]ComponentFactory

	RouteProcessor *RouteProcessorCommunicator
	Storage *commonStorage.Communicator
}

var GlobalComponentManager = ComponentManager{
	GeneralRegisterHandlers: []RegisterHandler{
		sendStorageOnRegister,
	},
	RegisterHandlers: map[string][]RegisterHandler{
		"storage": {
			sendStorageToComponentsOnRegister,
		},
	},
	Components: make([]commonCommunication.RemoteComponent, 0),
	factorySetters: map[string]ComponentFactory {
		"route-processor": routeProcessorFactory,
		"storage":         storageFactory,
	},
}

func (componentManager *ComponentManager) ConnectComponent(ctx context.Context, componentType string, endpoint *protoCommon.Endpoint) error {
	factory, ok := componentManager.factorySetters[componentType]
	if !ok {
		return fmt.Errorf("unknown component type %s: %w", componentType, ErrComponentNotFound)
	}

	comp, err := componentManager.createConnection(ctx, endpoint)
	if err != nil {
		return err
	}

	remoteComp, err := factory(ctx, componentManager, comp)
	if err != nil {
		return err
	}

	componentManager.Components = append(componentManager.Components, remoteComp)

	componentManager.runRegisterHandlers(ctx, componentType, remoteComp, endpoint)
	return nil
}

func (componentManager *ComponentManager) runRegisterHandlers(ctx context.Context, componentType string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
	for _, handler := range componentManager.GeneralRegisterHandlers {
		handler(ctx, componentManager, componentType, component, endpoint)
	}

	handlers, ok := componentManager.RegisterHandlers[componentType]
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(ctx, componentManager, componentType, component, endpoint)
	}
}

func (componentManager *ComponentManager) createConnection(ctx context.Context, endpoint *protoCommon.Endpoint) (*commonCommunication.ComponentCommunicator, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", endpoint.Host, endpoint.Port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not create connection to component: %w", err)
	}

	component := commonCommunication.NewComponentCommunicator(conn)
	err = component.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not ping component: %w", err)
	}

	return component, nil
}

func routeProcessorFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator) (commonCommunication.RemoteComponent, error) {
	manager.RouteProcessor =  NewRouteProcessorCommunicator(communicator)
	return manager.RouteProcessor, nil
}

func storageFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator) (commonCommunication.RemoteComponent, error) {
	manager.Storage = commonStorage.NewCommunicator(communicator)
	return manager.Storage, nil
}

func sendStorageOnRegister(ctx context.Context, manager *ComponentManager, componentType string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
	logger.Info("common handler")
}

func sendStorageToComponentsOnRegister(ctx context.Context, manager *ComponentManager, componentType string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
	logger.Info("storage")
}
