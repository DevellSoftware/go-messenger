package pubsub

import (
	"github.com/DevellSoftware/go-messenger/pkg/event"
	"github.com/DevellSoftware/go-messenger/pkg/messaging"
)

type PubSub struct {
	messagener  *messaging.Messenger
	subscribers []Subscriber
}

func NewPubSub(messenger *messaging.Messenger) *PubSub {
	return &PubSub{
		messagener:  messenger,
		subscribers: make([]Subscriber, 0),
	}
}

func (p *PubSub) Subscribe(subscriber Subscriber) {
	p.subscribers = append(p.subscribers, subscriber)
}

func (p *PubSub) Publish(event event.Event) {
	p.messagener.Send(event)
}
