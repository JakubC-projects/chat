package pubsub

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/JakubC-projects/chat/chat"
	"github.com/samber/lo"
)

type Pubsub struct {
	eventSource chat.EventSource

	subscribers   []*Subscriber
	subscribersMu *sync.Mutex
}

func New(eventSource chat.EventSource) *Pubsub {
	return &Pubsub{
		eventSource:   eventSource,
		subscribers:   []*Subscriber{},
		subscribersMu: &sync.Mutex{},
	}
}

var _ chat.Publisher = (*Pubsub)(nil)

type Subscriber struct {
	isClosed atomic.Bool
	topic    string
	channel  chan chat.Event
}

func (p *Pubsub) Subscribe(ctx context.Context, topic string) (chat.EventSource, error) {
	channel := make(chan chat.Event, 10)
	subscriber := &Subscriber{
		topic:   topic,
		channel: channel,
	}
	err := p.addSubscriber(subscriber)

	return subscriber, err
}

func (p *Pubsub) addSubscriber(sub *Subscriber) error {
	p.subscribersMu.Lock()
	defer p.subscribersMu.Unlock()
	p.subscribers = append(p.subscribers, sub)
	return nil

}

func (s *Subscriber) NextEvent(ctx context.Context) (chat.Event, error) {
	event := <-s.channel

	return event, nil
}

func (s *Subscriber) Close(ctx context.Context) {
	s.isClosed.Store(true)
}

func (p *Pubsub) Run() {
	for {
		event, err := p.eventSource.NextEvent(context.Background())
		if err != nil {
			continue
		}

		for _, sub := range p.subscribers {
			if sub.topic == event.Type && !sub.isClosed.Load() {
				sub.channel <- event
			}
		}
		p.cleanupSubscribers()
	}
}

func (p *Pubsub) cleanupSubscribers() {
	p.subscribersMu.Lock()
	defer p.subscribersMu.Unlock()
	p.subscribers = lo.Filter(p.subscribers, func(sub *Subscriber, _ int) bool {
		return !sub.isClosed.Load()
	})
}
