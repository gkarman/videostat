package dispatcher

import (
	"context"
	"reflect"
)

type Handler func(ctx context.Context, e any)

type Dispatcher struct {
	handlers map[reflect.Type][]Handler
}

func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[reflect.Type][]Handler),
	}
}

func normalize(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func (d *Dispatcher) Register(eventType any, h Handler) {
	t := normalize(reflect.TypeOf(eventType))
	d.handlers[t] = append(d.handlers[t], h)
}

func (d *Dispatcher) Dispatch(ctx context.Context, events []any) {
	for _, event := range events {
		t := normalize(reflect.TypeOf(event))
		if hs, ok := d.handlers[t]; ok {
			for _, h := range hs {
				h(ctx, event)
			}
		}
	}
}
