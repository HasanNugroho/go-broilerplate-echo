package modules

import (
	"sync"
)

type EventHandler func(payload any)

type EventBus struct {
	listeners map[string][]EventHandler
	lock      sync.RWMutex
}

func EventNew() *EventBus {
	return &EventBus{
		listeners: make(map[string][]EventHandler),
	}
}

func (bus *EventBus) On(event string, handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.listeners[event] = append(bus.listeners[event], handler)
}

func (bus *EventBus) Emit(event string, payload any) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	if handlers, ok := bus.listeners[event]; ok {
		for _, handler := range handlers {
			go handler(payload) // Async
		}
	}
}
