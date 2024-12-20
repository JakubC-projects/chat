package mock

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/chat/chat"
)

type EventSource struct {
	channel chan chat.Event
}

var _ chat.EventSource = (*EventSource)(nil)

func NewEventSource() *EventSource {
	return &EventSource{
		channel: make(chan chat.Event, 10),
	}
}

func (e *EventSource) NextEvent(ctx context.Context) (chat.Event, error) {
	evt, ok := <-e.channel
	if !ok {
		return evt, fmt.Errorf("source stopped")
	}
	return evt, nil
}

func (e *EventSource) Close(ctx context.Context) {
	close(e.channel)
}

func (e *EventSource) Send(evt chat.Event) {
	e.channel <- evt
}
