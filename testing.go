package eventbus

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertSubscriber(t *testing.T, b *BusMock, e Event, h Handler) {
	if !b.HasSubscriber(e, h) {
		msg := ""
		for s := range b.GetSubscribers(e) {
			msg += fmt.Sprintf(" %v", reflect.TypeOf(s))
		}
		t.Errorf("expected %v subscribed to %v, but only got:%s", reflect.TypeOf(h), reflect.TypeOf(e), msg)
	}
}

func AssertNotSubscriber(t *testing.T, b *BusMock, e Event, h Handler) {
	if b.HasSubscriber(e, h) {
		t.Errorf("asserted %v to be no subscriber of %v", reflect.TypeOf(h), reflect.TypeOf(e))
	}
}

func AssertNilSubscriber(t *testing.T, b *BusMock, e Event) {
	if b.GetSubscribers(e) != nil {
		msg := ""
		for s := range b.GetSubscribers(e) {
			msg += fmt.Sprintf(" %v", reflect.TypeOf(s))
		}
		t.Errorf("expected nil subscriber on %v, got:%s", reflect.TypeOf(e), msg)
	}
}

func AssertPublished(t *testing.T, b *BusMock, want []Event) {
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

func AssertEventsHandled(t *testing.T, h *HandlerMock, events ...Event) {
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

