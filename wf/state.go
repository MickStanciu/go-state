package wf

import "fmt"

type StateName string

const (
	stateInitial StateName = "INITIAL_STATE"
)

type State struct {
	name    StateName
	actions map[Event]*State
}

func (s *State) attachEvent(event Event, nextState *State) bool {
	_, ok := s.actions[event]
	if !ok {
		s.actions[event] = nextState
	}
	return !ok
}

func (s *State) execEvent(event Event) (*State, error) {
	newState, ok := s.actions[event]
	if ok {
		return newState, nil
	}
	return nil, fmt.Errorf("event %q is not defined for the current state %q", event, s.name)
}

func (s *State) GetName() StateName {
	return s.name
}
