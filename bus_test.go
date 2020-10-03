package eventbus

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

type PrimaryEvent struct{}
type SecondaryEvent struct{}

func (*PrimaryEvent) IsEvent()   {}
func (*SecondaryEvent) IsEvent() {}

// TestBus_Subscribe checks if a handler gets subscribed
// correctly. There should ever only be one instance of
// a certain handler for a certain event. It should not
// be possible to subscribe the same instance of a Handler
// twice.
func TestBus_Subscribe(t *testing.T) {
	h1 := &HandlerMock{}
	h2 := &HandlerMock{}
	e := &PrimaryEvent{}
	se := &SecondaryEvent{}
	b := NewMock()

	b.Subscribe(e, h1)
	// Should be subscribed.
	assertSubscriber(t, b, e, h1)

	b.Subscribe(e, h1)
	// Should not subscribe a second time.
	if len(b.GetSubscribers(e)) != 1 {
		t.Errorf("expected %d subscriber, got: %d", 1, len(b.GetSubscribers(e)))
	}
	// Should still be subscribed.
	assertSubscriber(t, b, e, h1)

	// Should not subscribe on different instance
	// of same event type.
	b.Subscribe(&PrimaryEvent{}, h1)
	// Should not subscribe a second time.
	if len(b.GetSubscribers(e)) != 1 {
		t.Errorf("expected %d subscriber, got: %d", 1, len(b.GetSubscribers(e)))
	}
	// Should still be subscribed.
	assertSubscriber(t, b, e, h1)

	// Subscribe a second handler to the same event.
	b.Subscribe(e, h2)
	// Should be subscribed.
	assertSubscriber(t, b, e, h2)

	// Subscribe same handler to different event.
	b.Subscribe(se, h1)
	// Should be subscribed.
	assertSubscriber(t, b, se, h1)
}

// TestBus_SubscribeConcurrent checks that the subscription
// is performed as expected under concurrent load.
func TestBus_SubscribeConcurrent(t *testing.T) {
	var NCHECKS = rand.Intn(10000)
	b := NewMock()
	e := &PrimaryEvent{}
	for i := 0; i < NCHECKS; i++ {
		h := &HandlerMock{}
		b.Subscribe(e, h)
	}

	if len(b.GetSubscribers(e)) != NCHECKS {
		t.Errorf("expected %d subscribers, got: %d", NCHECKS, len(b.GetSubscribers(e)))
	}
}

// TestBus_Unsubscribe makes sure that a handler gets
// correctly unsubscribed. There should be no error
// or panic if one tries to unsubscribe a unknown Handler
// or from a unknown Event.
func TestBus_Unsubscribe(t *testing.T) {
	h1 := &HandlerMock{}
	h2 := &HandlerMock{}
	e := &PrimaryEvent{}
	b := NewMock()

	// Subscribe two handler.
	b.Subscribe(e, h1)
	b.Subscribe(e, h2)

	// Unsubscribe first handler.
	b.Unsubscribe(e, h1)
	// h1 should be unsubscribed, h2 not.
	assertNotSubscriber(t, b, e, h1)
	assertSubscriber(t, b, e, h2)

	// Unsubscribe first handler, second time.
	b.Unsubscribe(e, h1)
	// h1 should be unsubscribed, h2 not.
	assertNotSubscriber(t, b, e, h1)
	assertSubscriber(t, b, e, h2)

	// Unsubscribe second handler.
	b.Unsubscribe(e, h2)
	// Should have no more subscribers.
	assertNilSubscriber(t, b, e)
}

func TestBus_Publish(t *testing.T) {
	events := []Event{&PrimaryEvent{}, &SecondaryEvent{}, &PrimaryEvent{}}

	wg := &sync.WaitGroup{}
	b := NewMock().WithWaitGroup(wg)
	h := &HandlerMock{wg: wg}

	// We need a second handler to make sure that the
	// waitGroup gets done.
	h2 := &HandlerMock{wg: wg}
	b.Subscribe(&SecondaryEvent{}, h2)

	// Assure that only subscribed events get published.
	b.Subscribe(&PrimaryEvent{}, h)
	for _, evt := range events {
		b.Publish(evt)
	}
	wg.Wait()
	assertPublished(t, b, events)
	assertEventsHandled(t, h, &PrimaryEvent{}, &PrimaryEvent{})

	// Assure that all Handler can be called on multiple
	// events.
	b = NewMock().WithWaitGroup(wg)
	h = &HandlerMock{wg: wg}
	b.Subscribe(&PrimaryEvent{}, h)
	b.Subscribe(&SecondaryEvent{}, h)
	for _, evt := range events {
		b.Publish(evt)
	}
	wg.Wait()
	assertEventsHandled(t, h, &PrimaryEvent{}, &PrimaryEvent{}, &SecondaryEvent{})
}

func assertSubscriber(t *testing.T, b *BusMock, e Event, h Handler) {
	if !b.HasSubscriber(e, h) {
		msg := ""
		for s := range b.GetSubscribers(e) {
			msg += fmt.Sprintf(" %v", reflect.TypeOf(s))
		}
		t.Errorf("expected %v subscribed to %v, but only got:%s", reflect.TypeOf(h), reflect.TypeOf(e), msg)
	}
}

func assertNotSubscriber(t *testing.T, b *BusMock, e Event, h Handler) {
	if b.HasSubscriber(e, h) {
		t.Errorf("asserted %v to be no subscriber of %v", reflect.TypeOf(h), reflect.TypeOf(e))
	}
}

func assertNilSubscriber(t *testing.T, b *BusMock, e Event) {
	if b.GetSubscribers(e) != nil {
		msg := ""
		for s := range b.GetSubscribers(e) {
			msg += fmt.Sprintf(" %v", reflect.TypeOf(s))
		}
		t.Errorf("expected nil subscriber on %v, got:%s", reflect.TypeOf(e), msg)
	}
}

func assertPublished(t *testing.T, b *BusMock, want []Event) {
	pe := b.GetPublished()
	msg := ""
	for i, e := range want {
		if reflect.TypeOf(e) != reflect.TypeOf(pe[i]) {
			msg += fmt.Sprintf("\tat pos %d want:%v, got:%v\n", i, reflect.TypeOf(e), reflect.TypeOf(pe[i]))
		}
	}
	if msg != "" {
		t.Errorf("error asserting published types:\n%s", msg)
	}
}

func assertEventsHandled(t *testing.T, h *HandlerMock, events ...Event) {
	// Make a copy of the handler events
	handledEvents := make([]Event, len(h.handled))
	copy(handledEvents, h.handled)

	eventHandled := make([]bool, len(events))
	for i, evt := range events {
		for j, hEvt := range handledEvents {
			if reflect.TypeOf(evt) == reflect.TypeOf(hEvt) {
				eventHandled[i] = true
				handledEvents = removeFromSlice(j, handledEvents)
				break // Event was found so we can break
			}
		}
	}

	var msg = ""
	var hasError = false
	for i, isHandled := range eventHandled {
		if !isHandled {
			msg += fmt.Sprintf("%v\n", reflect.TypeOf(events[i]))
			hasError = true
		}
	}

	if hasError {
		t.Errorf("events not handled: %s", msg)
	}
}

func removeFromSlice(pos int, slice []Event) []Event {
	slice[pos] = slice[len(slice)-1]
	slice[len(slice)-1] = nil
	slice = slice[:len(slice)-1]
	return slice
}
