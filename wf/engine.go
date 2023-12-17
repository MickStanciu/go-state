package wf

import "fmt"

type Engine struct {
	initialState *State
	currentState *State
	states       map[StateName]*State
}

type EngineOption func(state *Engine) error

const defaultInitialState = "STATE_INITIAL"

// NewEngine - build a new wf engine with a default initial state
func NewEngine(opts ...EngineOption) (*Engine, error) {
	initialState := &State{
		name:    defaultInitialState,
		events:  map[EventName]*State{},
		actions: map[EventName]func(StateName) error{},
	}

	e := &Engine{
		initialState: initialState,
		currentState: initialState,
		states: map[StateName]*State{
			defaultInitialState: initialState,
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

// WithInitialState - will build an initial state
func WithInitialState(initialStateName StateName) EngineOption {
	return func(e *Engine) error {
		initialState := e.getOrCreateState(initialStateName)
		initialState.events = map[EventName]*State{}
		e.initialState = initialState
		e.currentState = initialState
		return nil
	}
}

// WithState - will add a state during build
func WithState(fromStateName StateName, eventName EventName, toStateName StateName) EngineOption {
	return func(e *Engine) error {
		_, err := e.RegisterState(fromStateName, eventName, toStateName)
		return err
	}
}

// WithStateAndAction - will add a state & action during build
func WithStateAndAction(fromStateName StateName, eventName EventName, toStateName StateName, fn func(StateName) error) EngineOption {
	return func(e *Engine) error {
		if fn == nil {
			return fmt.Errorf("action cannot be nil")
		}

		_, err := e.RegisterState(fromStateName, eventName, toStateName)
		if err != nil {
			return err
		}

		fromState := e.GetState(fromStateName)
		if err := fromState.attachAction(eventName, fn); err != nil {
			return fmt.Errorf("an action is already defined for the state %q and action %q", fromStateName, eventName)
		}

		return err
	}
}

// WithExistingEngine - will replace the states / events / actions from another engine
func WithExistingEngine(ne *Engine) EngineOption {
	return func(e *Engine) error {
		//TODO? shall to a value copy instead of pointer copy ?
		e.states = ne.states
		e.currentState = ne.currentState
		e.initialState = ne.initialState
		return nil
	}
}

// RegisterState - will add a new state
// will return the new state or error if the state was previously defined
func (e *Engine) RegisterState(fromStateName StateName, eventName EventName, toStateName StateName) (*State, error) {
	fromState := e.getOrCreateState(fromStateName)
	toState := e.getOrCreateState(toStateName)
	if err := fromState.attachEvent(eventName, toState); err != nil {
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
		events:  map[EventName]*State{},
		actions: map[EventName]func(StateName) error{},
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
	nextState, err := e.currentState.execEvent(event)
	if err != nil {
		return nil, err
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

// AttachAction - will attach an action to existing state/event
func (e *Engine) AttachAction(eventName EventName, fn func(nextStateName StateName) error) error {
	if err := e.currentState.attachAction(eventName, fn); err != nil {
		return fmt.Errorf("an action is already defined for the state %q and action %q", e.currentState.GetName(), eventName)
	}
	return nil
}
