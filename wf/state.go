package wf

import "fmt"

type StateName string

//type EventAction func() error

type State struct {
	name    StateName
	events  map[EventName]*State
	actions map[EventName]func() error
}

func (s *State) attachEvent(event EventName, nextState *State) error {
	_, ok := s.events[event]
	if !ok {
		s.events[event] = nextState
		return nil
	}
	return fmt.Errorf("event %q is already defined for the state %q", event, s.name)
}

func (s *State) attachAction(event EventName, fn func() error) error {
	_, ok := s.actions[event]
	if ok {
		return fmt.Errorf("action is already defined for the event %q and the state %q", event, s.name)
	}

	s.actions[event] = fn
	return nil
}

func (s *State) execEvent(event EventName) (*State, error) {
	newState, ok := s.events[event]
	if !ok {
		return nil, fmt.Errorf("event %q is not defined for the state %q", event, s.name)
	}

	actionFn, ok := s.actions[event]
	if !ok {
		return newState, nil
	}

	// exec action
	if err := actionFn(); err != nil {
		return nil, err
	}

	return newState, nil
}

func (s *State) GetName() StateName {
	return s.name
}
