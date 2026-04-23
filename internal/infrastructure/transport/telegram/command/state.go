package command

import "sync"

type userState struct {
	PlatformName string
	WaitingURL   bool
	WaitingVideoURL bool
}

type State struct {
	mu     sync.Mutex
	states map[int64]*userState
}

func NewState() *State {
	return &State{
		states: make(map[int64]*userState),
	}
}

func (s *State) Set(chatID int64, st *userState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[chatID] = st
}

func (s *State) Get(chatID int64) (*userState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[chatID]
	return st, ok
}

func (s *State) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, chatID)
}