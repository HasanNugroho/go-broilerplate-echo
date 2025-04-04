package modules

// Event represents a generic event structure
type Event struct {
	Type string
	Data interface{}
}

// EventBus represents the event bus that handles event subscription and dispatching
type EventBus struct {
	subscribers map[string][]chan<- Event
}

// NewEventBus creates a new instance of the event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
	}
}

// Subscribe adds a new subscriber for a given event type
func (eb *EventBus) Subscribe(eventType string, subscriber chan<- Event) {
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

// Publish sends an event to all subscribers of a given event type
func (eb *EventBus) Publish(event Event) {
	subscribers := eb.subscribers[event.Type]
	for _, subscriber := range subscribers {
		subscriber <- event
	}
}
