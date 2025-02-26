package event

import (
	"context"
	"log"
	"sync"
)

type Name string

type Event struct {
	Name Name
	Data interface{}
}

type Listener interface {
	Listen(ctx context.Context, event Event)
}

type Dispatcher struct {
	jobs   chan Event
	events map[Name][]Listener
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewDispatcher(workerCount int) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	d := &Dispatcher{
		jobs:   make(chan Event, 100), // Buffered channel
		events: make(map[Name][]Listener),
		ctx:    ctx,
		cancel: cancel,
	}

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go d.worker()
	}

	return d
}

func (d *Dispatcher) Register(listener Listener, names ...Name) {
	for _, name := range names {
		d.events[name] = append(d.events[name], listener)
	}
}

func (d *Dispatcher) Dispatch(event Event) {
	d.jobs <- event
}

func (d *Dispatcher) worker() {
	for {
		select {
		case event := <-d.jobs:
			if listeners, exists := d.events[event.Name]; exists {
				for _, listener := range listeners {
					d.wg.Add(1)
					go func(listener Listener, event Event) {
						defer d.wg.Done()
						defer func() {
							if r := recover(); r != nil {
								log.Printf("Recovered from panic in listener: %v", r)
							}
						}()
						listener.Listen(d.ctx, event)
					}(listener, event)
				}
			}
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) Shutdown() {
	d.cancel()
	close(d.jobs)
	d.wg.Wait()
}
