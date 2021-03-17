package communication

import (
	"context"
	"errors"
	"fmt"
	commonCommunication "github.com/kulycloud/common/communication"
	protoCommon "github.com/kulycloud/protocol/common"
	"google.golang.org/grpc"
)

var ErrComponentNotFound = errors.New("component not found")

type RegisterHandler = func(context.Context, *ComponentManager, string, commonCommunication.RemoteComponent, *protoCommon.Endpoint)
type ComponentFactory = func(context.Context, *ComponentManager, *commonCommunication.ComponentCommunicator, *protoCommon.Endpoint) (commonCommunication.RemoteComponent, error)

type ComponentManager struct {
	GeneralRegisterHandlers []RegisterHandler
	RegisterHandlers        map[string][]RegisterHandler
	Components              []commonCommunication.RemoteComponent
	factorySetters          map[string]ComponentFactory

	ApiServer        *ApiServerCommunicator
	Storage          *commonCommunication.StorageCommunicator
	storageEndpoints []*protoCommon.Endpoint
	Ingress          *IngressCommunicator
	ServiceManager   *ServiceManagerCommunicator
}

var GlobalComponentManager = ComponentManager{
	GeneralRegisterHandlers: []RegisterHandler{
		sendStorageOnRegister,
	},
	RegisterHandlers: map[string][]RegisterHandler{
		"storage": {
			sendStorageToComponentsOnRegister,
		},
		"service-manager": {
			func(ctx context.Context, manager *ComponentManager, typeName string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
				manager.ServiceManager.ReconcileNamespace(ctx, "core")
				manager.ServiceManager.ReconcileNamespace(ctx, "u01")
			},
		},
	},
	Components: make([]commonCommunication.RemoteComponent, 0),
	factorySetters: map[string]ComponentFactory{
		"api-server":      apiServerFactory,
		"storage":         storageFactory,
		"ingress":         ingressFactory,
		"service-manager": serviceManagerFactory,
	},
	Storage: commonCommunication.NewEmptyStorageCommunicator(),
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

	remoteComp, err := factory(ctx, componentManager, comp, endpoint)
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

func apiServerFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator, endpoint *protoCommon.Endpoint) (commonCommunication.RemoteComponent, error) {
	manager.ApiServer = NewApiServerCommunicator(communicator)
	return manager.ApiServer, nil
}

func storageFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator, endpoint *protoCommon.Endpoint) (commonCommunication.RemoteComponent, error) {
	manager.Storage.UpdateComponentCommunicator(communicator)
	manager.storageEndpoints = []*protoCommon.Endpoint{endpoint}
	return manager.Storage, nil
}

func ingressFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator, endpoint *protoCommon.Endpoint) (commonCommunication.RemoteComponent, error) {
	manager.Ingress = NewIngressCommunicator(communicator)
	return manager.Ingress, nil
}

func serviceManagerFactory(ctx context.Context, manager *ComponentManager, communicator *commonCommunication.ComponentCommunicator, endpoint *protoCommon.Endpoint) (commonCommunication.RemoteComponent, error) {
	manager.ServiceManager = NewServiceManagerCommunicator(communicator)
	return manager.ServiceManager, nil
}

func sendStorageOnRegister(ctx context.Context, manager *ComponentManager, componentType string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
	// Send storage to component that was just registered when storage already available (except it is a storage)
	if componentType != "storage" && manager.Storage.Ready() {
		err := component.RegisterStorageEndpoints(ctx, manager.storageEndpoints)
		if err != nil {
			logger.Warnw("Could not propagate storage to endpoint", "componentType", componentType, "endpoint", endpoint, "error", err)
		}
	}
}

func sendStorageToComponentsOnRegister(ctx context.Context, manager *ComponentManager, componentType string, component commonCommunication.RemoteComponent, endpoint *protoCommon.Endpoint) {
	// When storage connects propagate it to all cluster components
	if manager.Storage.Ready() {
		for _, component := range manager.Components {
			err := component.RegisterStorageEndpoints(ctx, manager.storageEndpoints)
			if err != nil {
				logger.Warnw("Could not propagate storage to endpoint", "componentType", componentType, "endpoint", endpoint, "error", err)
			}
		}
	}
	logger.Info("storage")
}
