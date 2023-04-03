package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type MessengerCallback func(ctx goka.Context, msg interface{})

type Messenger struct {
	callback  MessengerCallback
	emitter   *goka.Emitter
	processor *goka.Processor
	cancel    context.CancelFunc
	context   context.Context
}

func NewMessenger(brokerHost string, topic string, group string, callback MessengerCallback) *Messenger {
	m := &Messenger{}
	m.callback = callback

	gokaCallback := func(ctx goka.Context, msg interface{}) {
		m.callback(ctx, msg)
	}

	var gokaGroup goka.Group = goka.Group(fmt.Sprint(group))
	var gokaTopic goka.Stream = goka.Stream(fmt.Sprint(topic))

	g := goka.DefineGroup(gokaGroup,
		goka.Input(gokaTopic, new(codec.String), gokaCallback),
	)

	p, err := goka.NewProcessor([]string{brokerHost}, g)

	m.processor = p

	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	m.context = ctx
	m.cancel = cancel

	/*
		wait := make(chan os.Signal, 1)
		signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
		<-wait   // wait for SIGINT/SIGTERM
		cancel() // gracefully stop processor
		<-done
	*/
	emitter, err := goka.NewEmitter([]string{brokerHost}, gokaTopic, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	m.emitter = emitter

	return m
}

func (m *Messenger) Send(message interface{}) {
	err := m.emitter.EmitSync("some-key", message)

	if err != nil {
		log.Fatalf("error emitting message: %v", err)
	}
}

func (m *Messenger) Listen() {
	done := make(chan struct{})

	go func() {
		defer close(done)

		if err := m.processor.Run(m.context); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	<-done
}

func (m *Messenger) Close() {
	m.cancel()
	m.emitter.Finish()
}
