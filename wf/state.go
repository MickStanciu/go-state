package wf

type StateName string

//type EventAction func() error

type State struct {
	name   StateName
	events map[EventName]*State
	//actions map[EventName] EventAction
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
