package pubsub

import "github.com/DevellSoftware/go-messenger/pkg/event"

type Subscriber struct {
	name     string
	callback func(event event.Event)
}
