package eventbus

import (
	"sync"
)

type BusMock struct {
	sync.RWMutex
	published         []Event
	subscribeCalled   int
	unsubscribeCalled int
	bus               *Bus

	wg *sync.WaitGroup
}

func NewMock() *BusMock {
	return &BusMock{bus: New()}
}

func (b *BusMock) WithWaitGroup(wg *sync.WaitGroup) *BusMock {
	b.wg = wg
	return b
}

func (b *BusMock) Subscribe(e interface{ Event }, sub interface{ Handler }) {
	b.subscribeCalled++
	b.bus.Subscribe(e, sub)
}

func (b *BusMock) Unsubscribe(e interface{ Event }, sub interface{ Handler }) {
	b.unsubscribeCalled++
	b.bus.Unsubscribe(e, sub)
}

func (b *BusMock) Publish(e interface{ Event }) {
	b.Lock()
	defer b.Unlock()
	if b.wg != nil {
		b.wg.Add(1)
	}
	b.published = append(b.published, e)
	b.bus.Publish(e)
}

func (b *BusMock) GetSubscribers(e interface{ Event }) map[interface{ Handler }]struct{} {
	b.RLock()
	defer b.RUnlock()
	return b.bus.subscribers[e]
}

func (b *BusMock) GetPublished() []Event {
	return b.published
}

func (b *BusMock) HasSubscriber(event interface{ Event }, sub interface{ Handler }) bool {
	subs := b.GetSubscribers(event)
	for s := range subs {
		if s == sub {
			return true
		}
	}
	return false
}

type HandlerMock struct {
	sync.Mutex
	handled []Event

	wg *sync.WaitGroup
}

func (h *HandlerMock) Handle(e interface{Event}) {
	h.Lock()
	defer h.Unlock()
	h.handled = append(h.handled, e)
	if h.wg != nil {
		h.wg.Done()
	}
}
