package eventbus

import (
	"math/rand"
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
	AssertSubscriber(t, b, e, h1)

	b.Subscribe(e, h1)
	// Should not subscribe a second time.
	if len(b.GetSubscribers(e)) != 1 {
		t.Errorf("expected %d subscriber, got: %d", 1, len(b.GetSubscribers(e)))
	}
	// Should still be subscribed.
	AssertSubscriber(t, b, e, h1)

	// Should not subscribe on different instance
	// of same event type.
	b.Subscribe(&PrimaryEvent{}, h1)
	// Should not subscribe a second time.
	if len(b.GetSubscribers(e)) != 1 {
		t.Errorf("expected %d subscriber, got: %d", 1, len(b.GetSubscribers(e)))
	}
	// Should still be subscribed.
	AssertSubscriber(t, b, e, h1)

	// Subscribe a second handler to the same event.
	b.Subscribe(e, h2)
	// Should be subscribed.
	AssertSubscriber(t, b, e, h2)

	// Subscribe same handler to different event.
	b.Subscribe(se, h1)
	// Should be subscribed.
	AssertSubscriber(t, b, se, h1)
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
	AssertNotSubscriber(t, b, e, h1)
	AssertSubscriber(t, b, e, h2)

	// Unsubscribe first handler, second time.
	b.Unsubscribe(e, h1)
	// h1 should be unsubscribed, h2 not.
	AssertNotSubscriber(t, b, e, h1)
	AssertSubscriber(t, b, e, h2)

	// Unsubscribe second handler.
	b.Unsubscribe(e, h2)
	// Should have no more subscribers.
	AssertNilSubscriber(t, b, e)
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
	AssertPublished(t, b, events)
	AssertEventsHandled(t, h, &PrimaryEvent{}, &PrimaryEvent{})

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
	AssertEventsHandled(t, h, &PrimaryEvent{}, &PrimaryEvent{}, &SecondaryEvent{})
}
