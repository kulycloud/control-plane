package communication

import (
	"context"
	"fmt"
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

func (communicator *ServiceManagerCommunicator) ReconcileNamespace(ctx context.Context, namespace string) error {
	_, err := communicator.client.Reconcile(ctx, &protoServices.ReconcileRequest{
		Namespace: namespace,
	})

	if err != nil {
		return fmt.Errorf("error from service-manager: %w", err)
	}
	return nil
}
