package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/JakubC-projects/chat/chat"
	"github.com/JakubC-projects/chat/chat/mock"
	"github.com/stretchr/testify/assert"
)

func TestPubsub(t *testing.T) {
	eventSrc := mock.NewEventSource()

	pubsub := New(eventSrc)
	go pubsub.Run()

	sub1, err := pubsub.Subscribe(context.Background(), "test")
	events1 := []chat.Event{}
	assert.NoError(t, err)
	go func() {
		for {
			evt, err := sub1.NextEvent(context.Background())
			assert.NoError(t, err)
			events1 = append(events1, evt)
		}
	}()

	sub2, err := pubsub.Subscribe(context.Background(), "test")
	events2 := []chat.Event{}
	assert.NoError(t, err)
	go func() {
		for {
			evt, err := sub2.NextEvent(context.Background())
			assert.NoError(t, err)
			events2 = append(events2, evt)
		}
	}()

	events := []chat.Event{
		{Type: "test", Payload: "Hello world"},
		{Type: "test", Payload: "Hello world2"},
	}
	for _, e := range events {
		eventSrc.Send(e)
	}
	eventSrc.Send(chat.Event{Type: "other", Payload: "Hello world3"})

	time.Sleep(1 * time.Second)
	assert.Equal(t, events, events1)
	assert.Equal(t, events, events2)
}
