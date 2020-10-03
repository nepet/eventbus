package eventbus

import (
	"fmt"
	"reflect"
)

// This is an example for a handler that uses a channel
// to pass an event out of the handler struct.
type chanHandler struct {
	evt chan<- Event
}

func (c *chanHandler) Handle(event interface{ Event }) {
	c.evt <- event
}

type exampleEvent struct{}

func (*exampleEvent) IsEvent() {}

func ExampleHandler() {
	bus := New()

	// Create an instance of a handler
	eChan := make(chan Event)
	handler := &chanHandler{evt: eChan}

	bus.Subscribe(&exampleEvent{}, handler)
	defer bus.Unsubscribe(&exampleEvent{}, handler)

	bus.Publish(&exampleEvent{})

	evt := <-eChan
	fmt.Println(reflect.TypeOf(evt))
	// Output: *eventbus.exampleEvent
}
