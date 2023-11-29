package wf

import "fmt"

type Engine struct {
	initialState *State
	currentState *State
	states       map[StateName]*State
}

type EngineOption func(state *Engine) error

// NewEngine - build a new wf engine with an initial state
func NewEngine(initialStateName StateName, opts ...EngineOption) (*Engine, error) {
	initialState := &State{
		name:    initialStateName,
		actions: map[EventName]*State{},
	}

	e := &Engine{
		initialState: initialState,
		currentState: initialState,
		states: map[StateName]*State{
			initialStateName: initialState,
		},
	}

	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

// WithState - will add a state during build
func WithState(fromStateName StateName, eventName EventName, toStateName StateName) EngineOption {
	return func(e *Engine) error {
		_, err := e.RegisterState(fromStateName, eventName, toStateName)
		return err
	}
}

// RegisterState - will add a new state
// will return the new state or error if the state was previously defined
func (e *Engine) RegisterState(fromStateName StateName, eventName EventName, toStateName StateName) (*State, error) {
	fromState := e.getOrCreateState(fromStateName)
	toState := e.getOrCreateState(toStateName)
	if eventOk := fromState.attachEvent(eventName, toState); !eventOk {
		return nil, fmt.Errorf("event %q already defined for the state %q", eventName, fromState.name)
	}
	return toState, nil
}

// getOrCreateState - gets or creates a state.
func (e *Engine) getOrCreateState(name StateName) *State {
	state, ok := e.states[name]
	if ok {
		return state
	}

	newState := &State{
		name:    name,
		actions: map[EventName]*State{},
	}
	e.states[name] = newState
	return newState
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
	return e.initialState
}

// GetCurrentState - returns the current state
func (e *Engine) GetCurrentState() *State {
	return e.currentState
}

// ProcessEvent - will run the event and return the next state
// in case of the event was not defined, will return error
func (e *Engine) ProcessEvent(event EventName) (*State, error) {
	nextState, ok := e.currentState.execEvent(event)
	if !ok {
		return nil, fmt.Errorf("event %q is not defined for the current state %q", event, e.currentState.name)
	}
	e.currentState = nextState
	return nextState, nil
}

// JumpToState - will force a state change
func (e *Engine) JumpToState(name StateName) error {
	s, ok := e.states[name]
	if !ok {
		return fmt.Errorf("state %q is not defined", name)
	}
	e.currentState = s
	return nil
}
