package wf

import "fmt"

type Engine struct {
	currentState *State
	states       map[StateName]*State
}

// NewEngine - build a new wf engine with an initial state
func NewEngine() *Engine {
	s := map[StateName]*State{}

	state := &State{
		name:    stateInitial,
		actions: map[Event]*State{},
	}

	s[stateInitial] = state
	return &Engine{
		currentState: state,
		states:       s,
	}
}

// RegisterState - will add a new state
// will return the new state or error if the state was previously defined
func (e *Engine) RegisterState(name StateName) (*State, error) {
	_, ok := e.states[name]
	if ok {
		return nil, fmt.Errorf("state %q already defined", name)
	}

	s := &State{
		name:    name,
		actions: map[Event]*State{},
	}
	e.states[name] = s

	return s, nil
}

// RegisterEvent - will add an event to facilitate transition from current state to the next state
func (e *Engine) RegisterEvent(curState *State, event Event, nextState *State) error {
	if !curState.attachEvent(event, nextState) {
		return fmt.Errorf("event %q already defined for the state %q", event, curState.name)
	}
	return nil
}

// GetState - returns a state by name or nil
func (e *Engine) GetState(name StateName) *State {
	s, ok := e.states[name]
	if ok {
		return s
	}
	return nil
}

// GetInitialState - returns the initial state
func (e *Engine) GetInitialState() *State {
	return e.GetState(stateInitial)
}

// GetCurrentState - returns the current state
func (e *Engine) GetCurrentState() *State {
	return e.currentState
}

// ProcessEvent - will run the event and return the next state
// in case of the event was not defined, will return error
func (e *Engine) ProcessEvent(event Event) (*State, error) {
	nextState, err := e.currentState.execEvent(event)
	if err != nil {
		return nil, err
	}
	e.currentState = nextState
	return nextState, nil
}
