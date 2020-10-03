package eventbus

import (
	"fmt"
)

// This is an example of an event that contains some
// information in the struct
type dataEvent struct {
	someString string
	someInt    int
}

func (*dataEvent) IsEvent() {}

type printHandler struct{}

func (*printHandler) Handle(evt interface{ Event }) {
	if e, ok := evt.(*dataEvent); ok {
		fmt.Printf("got event with data: %s, %d\n", e.someString, e.someInt)
	}

}

func ExampleEvent() {
	bus := New()

	handler := &printHandler{}

	bus.Subscribe(&dataEvent{}, handler)
	defer bus.Unsubscribe(&dataEvent{}, handler)

	bus.Publish(&dataEvent{someString: "Data string", someInt: 12})
	// Output:
}
