package eventbus

import "sync"

type Event interface {
	IsEvent()
}

// Handler handles events published by the event bus.
type Handler interface {
	Handle(interface{ Event })
}

// Bus holds a mapping of the subscribers to event types. Bus
// also manages the un-/subscription of subscribers and the
// publishing of events.
type Bus struct {
	sync.RWMutex
	subscribers map[interface{ Event }]map[interface{ Handler }]struct{}
}

// New returns a new Bus.
func New() *Bus {
	return &Bus{
		subscribers: make(map[interface{ Event }]map[interface{ Handler }]struct{}),
	}
}

// Subscribe a subscriber to an event type. It is ensured that a handler
// can only be subscribed to a certain event once.
func (b *Bus) Subscribe(e interface{ Event }, sub interface{ Handler }) {
	b.Lock()
	defer b.Unlock()

	if _, exists := b.subscribers[e]; !exists {
		b.subscribers[e] = make(map[interface{ Handler }]struct{})
	}
	b.subscribers[e][sub] = struct{}{}
}

// Unsubscribe a subscriber from an event type.
func (b *Bus) Unsubscribe(e interface{ Event }, sub interface{ Handler }) {
	b.Lock()
	defer b.Unlock()

	if subs, exists := b.subscribers[e]; exists {
		delete(subs, sub)
		if len(subs) == 0 {
			delete(b.subscribers, e)
		}
	}
}

// Publish an event to all subscribed subscribers. This
// uses the Handler of an subscriber.
func (b *Bus) Publish(e interface{ Event }) {
	b.RLock()
	defer b.RUnlock()

	if subs, exists := b.subscribers[e]; exists {
		for sub, _ := range subs {
			go sub.Handle(e)
		}
	}
}
