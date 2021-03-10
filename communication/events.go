package communication

import (
	"fmt"

    protoCommon "github.com/kulycloud/protocol/common"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
)

type eventStream struct {
    destination string
    stream protoControlPlane.ControlPlane_RegisterComponentServer
    activelyListening map[string]bool
}

func newEventStream(endpoint *protoCommon.Endpoint, stream protoControlPlane.ControlPlane_RegisterComponentServer) *eventStream {
    return &eventStream{
        destination: fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port),
        stream: stream,
    }
}

func (es *eventStream) listenOnEvent(eventType string) {
    es.activelyListening[eventType] = true
}

func (es *eventStream) send(event *protoCommon.Event) error {
    return es.stream.Send(event)
}

