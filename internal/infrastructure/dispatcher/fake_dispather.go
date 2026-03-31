package dispatcher

import "context"

type FakeDispatcher struct {
	Events [][]any
}

func NewFakeDispatcher() *FakeDispatcher {
	return &FakeDispatcher{
		Events: make([][]any, 0),
	}
}

func (s *FakeDispatcher) Dispatch(ctx context.Context, events []any) {
	s.Events = append(s.Events, events)
}