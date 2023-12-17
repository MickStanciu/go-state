package wf_test

import (
	"fmt"
	"testing"

	"github.com/MickStanciu/go-state/wf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	START  = "START"
	S1     = "STATE_1"
	S2     = "STATE_2"
	S3     = "STATE_3"
	S4     = "STATE_4"
	FINISH = "FINISH"

	E1 = "EVENT_1"
	E2 = "EVENT_2"
	E3 = "EVENT_3"
	E4 = "EVENT_4"
	E5 = "EVENT_5"
	E6 = "EVENT_6"
)

func TestNewEngine(t *testing.T) {
	wfe, err := wf.NewEngine()
	require.NoError(t, err)
	require.NotNil(t, wfe)

	initial := wfe.GetInitialState()
	assert.NotNil(t, initial)
	assert.EqualValues(t, "STATE_INITIAL", initial.GetName())
}

func TestNewEngine_WithState_WhenNoErrors(t *testing.T) {
	wfe, err := wf.NewEngine(
		wf.WithInitialState("START"),
		wf.WithState(START, E1, S1),
		wf.WithState(S1, E2, S2),
		wf.WithState(S1, E4, S4),
		wf.WithState(S2, E3, S3),
		wf.WithState(S3, E3, S2),
		wf.WithState(S3, E5, FINISH),
		wf.WithState(S4, E5, FINISH),
	)
	require.NoError(t, err)
	require.NotNil(t, wfe)

	// test initial state
	initial := wfe.GetInitialState()
	assert.NotNil(t, initial)
	assert.EqualValues(t, "START", initial.GetName())

	// test jump to S1
	nextState, err := wfe.ProcessEvent(E1)
	require.NoError(t, err)
	require.NotNil(t, nextState)
	assert.EqualValues(t, wf.StateName("STATE_1"), nextState.GetName())

	// test jump to S2
	nextState, err = wfe.ProcessEvent(E2)
	require.NoError(t, err)
	require.NotNil(t, nextState)
	assert.EqualValues(t, wf.StateName("STATE_2"), nextState.GetName())

	// test illegal jump
	nextState, err = wfe.ProcessEvent(E2)
	require.Nil(t, nextState)
	require.NotNil(t, err)
	assert.EqualValues(t, "event \"EVENT_2\" is not defined for the state \"STATE_2\"", err.Error())
}

func TestEngine_RegisterState(t *testing.T) {
	wfe, err := wf.NewEngine(wf.WithInitialState("START"))
	require.NoError(t, err)
	require.NotNil(t, wfe)

	s0 := wfe.GetInitialState()
	s1, err := wfe.RegisterState(START, E1, S1)
	require.NoError(t, err)
	require.NotNil(t, s1)

	s1_ := wfe.GetState(S1)
	require.NotNil(t, s1_)

	assert.EqualValues(t, s0, wfe.GetCurrentState())
	assert.EqualValues(t, wf.StateName("STATE_1"), s1_.GetName())
}

func TestEngine_GetCurrentState(t *testing.T) {
	wfe, err := wf.NewEngine(wf.WithInitialState("START"))
	require.NoError(t, err)
	require.NotNil(t, wfe)

	s0 := wfe.GetInitialState()
	assert.NotNil(t, s0)
	assert.EqualValues(t, s0, wfe.GetCurrentState())

	s1, _ := wfe.RegisterState(START, E1, S1)
	assert.NotNil(t, s1)

	nextState, err := wfe.ProcessEvent(E1)
	require.NoError(t, err)
	assert.EqualValues(t, nextState, wfe.GetCurrentState())
}

func TestEngine_GetState_When_Error(t *testing.T) {
	wfe, err := wf.NewEngine(wf.WithInitialState("START"))
	require.NoError(t, err)
	require.NotNil(t, wfe)

	s1 := wfe.GetState("FAKE STATE")
	assert.Nil(t, s1)
}

func TestEngine_ProcessEvent_When_Error(t *testing.T) {
	wfe, err := wf.NewEngine(wf.WithInitialState("START"))
	require.NoError(t, err)
	require.NotNil(t, wfe)

	s1, _ := wfe.RegisterState(START, E1, S1)
	assert.NotNil(t, s1)

	_, err = wfe.ProcessEvent(E2)
	require.Error(t, err)
	assert.EqualValues(t, "event \"EVENT_2\" is not defined for the state \"START\"", err.Error())
}

func TestEngine_JumpToState(t *testing.T) {
	wfe, err := wf.NewEngine(wf.WithInitialState("START"))
	require.NoError(t, err)
	require.NotNil(t, wfe)

	wfe.GetInitialState()
	_, err = wfe.RegisterState(START, E1, S1)
	require.NoError(t, err)
	_, err = wfe.RegisterState(S1, E2, S2)
	require.NoError(t, err)

	err = wfe.JumpToState(S2)
	require.NoError(t, err)
	assert.EqualValues(t, "STATE_2", wfe.GetCurrentState().GetName())
}

func TestNewEngine_WithStateAndAction(t *testing.T) {
	val := 10
	pVal := &val

	doer := func(wf.StateName) error {
		*pVal++
		return nil
	}

	wfe, err := wf.NewEngine(
		wf.WithInitialState("START"),
		wf.WithStateAndAction(START, E1, S1, doer),
	)
	require.NoError(t, err)
	require.NotNil(t, wfe)

	assert.EqualValues(t, 10, val)
	_, err = wfe.ProcessEvent(E1)
	assert.EqualValues(t, 11, val)
}

func TestNewEngine_WithStateAndAction_ShouldFailWhenActionFails(t *testing.T) {
	doer := func(wf.StateName) error {
		return fmt.Errorf("some error")
	}

	wfe, err := wf.NewEngine(
		wf.WithInitialState("START"),
		wf.WithStateAndAction(START, E1, S1, doer),
	)
	require.NoError(t, err)
	require.NotNil(t, wfe)

	_, err = wfe.ProcessEvent(E1)
	assert.EqualValues(t, "some error", err.Error())

}
