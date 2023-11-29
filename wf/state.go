package wf

type StateName string

const (
	stateInitial StateName = "INITIAL_STATE"
)

type State struct {
	name    StateName
	actions map[EventName]*State
}

func (s *State) attachEvent(event EventName, nextState *State) bool {
	_, ok := s.actions[event]
	if !ok {
		s.actions[event] = nextState
	}
	return !ok
}

func (s *State) execEvent(event EventName) (*State, bool) {
	newState, ok := s.actions[event]
	if ok {
		return newState, true
	}
	return nil, false
}

func (s *State) GetName() StateName {
	return s.name
}
