package wf

type StateName string

type State struct {
	name   StateName
	events map[EventName]*State
}

func (s *State) attachEvent(event EventName, nextState *State) bool {
	_, ok := s.events[event]
	if !ok {
		s.events[event] = nextState
	}
	return !ok
}

func (s *State) execEvent(event EventName) (*State, bool) {
	newState, ok := s.events[event]
	if ok {
		return newState, true
	}
	return nil, false
}

func (s *State) GetName() StateName {
	return s.name
}
