package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"

	protoCommon "github.com/kulycloud/protocol/common"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
)

type eventStream struct {
	destination       string
	stream            protoControlPlane.ControlPlane_RegisterComponentServer
	activelyListening map[string]bool
}

func newEventStream(endpoint *protoCommon.Endpoint, stream protoControlPlane.ControlPlane_RegisterComponentServer) *eventStream {
	return &eventStream{
		destination:       commonCommunication.NewIdentifierFromEndpoint(endpoint),
		stream:            stream,
		activelyListening: map[string]bool{},
	}
}

func (es *eventStream) listenOnEvent(eventType string) {
	es.activelyListening[eventType] = true
}

func (es *eventStream) send(event *protoCommon.Event) error {
	if es.activelyListening[event.Type] {
		return es.stream.Send(event)
	}
	return nil
}

func (es *eventStream) sendConfirmation() error {
	return es.stream.Send(&protoCommon.Event{Type: "confirmation"})
}
