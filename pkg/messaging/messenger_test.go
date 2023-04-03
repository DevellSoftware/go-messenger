package messaging

import (
	"testing"

	"github.com/lovoo/goka"
)

func TestSendAndReceive(t *testing.T) {
	m1 := NewMessenger("localhost:9092", "test", "test", func(ctx goka.Context, msg interface{}) {
		t.Log("Received message: ", msg)
	})

	m2 := NewMessenger("localhost:9092", "test", "test", func(ctx goka.Context, msg interface{}) {
		t.Log("Received message: ", msg)
	})

	m1.Send(&Message{
		Topic:   "test",
		Payload: "Hello World",
	})
}
