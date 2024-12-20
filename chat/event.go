package chat

import "context"

type Event struct {
	Type    string
	Payload string
}

type EventSource interface {
	NextEvent(context.Context) (Event, error)
	Close(context.Context)
}

type Publisher interface {
	Subscribe(ctx context.Context, topic string) (EventSource, error)
}
