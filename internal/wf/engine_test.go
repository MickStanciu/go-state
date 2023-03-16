package wf_test

import (
	"testing"

	"go-state/internal/wf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	S1 = "STATE_1"
	S2 = "STATE_2"
	S3 = "STATE_3"
	S4 = "STATE_4"
	S5 = "STATE_5"

	E1 = "EVENT_1"
	E2 = "EVENT_2"
	E3 = "EVENT_3"
	E4 = "EVENT_4"
	E5 = "EVENT_5"
	E6 = "EVENT_6"
)

func TestNewEngine(t *testing.T) {
	wfe := wf.NewEngine()
	require.NotNil(t, wfe)

	s0 := wfe.GetInitialState()
	assert.NotNil(t, s0)
}

func TestEngine_RegisterState(t *testing.T) {
	wfe := wf.NewEngine()
	s0 := wfe.GetInitialState()
	s1, err := wfe.RegisterState(S1)
	require.NoError(t, err)

	require.NotNil(t, s1)
	assert.EqualValues(t, s0, wfe.GetCurrentState())
	assert.EqualValues(t, wf.StateName("STATE_1"), s1.GetName())
}

func TestEngine_RegisterState_When_Error(t *testing.T) {
	wfe := wf.NewEngine()
	s0 := wfe.GetInitialState()
	_, err := wfe.RegisterState(s0.GetName())
	require.Error(t, err)
	assert.EqualValues(t, "state \"INITIAL_STATE\" already defined", err.Error())
}

func TestEngine_RegisterEvent_When_Error(t *testing.T) {
	wfe := wf.NewEngine()

	s0 := wfe.GetInitialState()
	s1, _ := wfe.RegisterState(S1)
	err := wfe.RegisterEvent(s0, E1, s1)
	require.NoError(t, err)
	err = wfe.RegisterEvent(s0, E1, s1)
	assert.Error(t, err)
	assert.EqualValues(t, "event \"EVENT_1\" already defined for the state \"INITIAL_STATE\"", err.Error())
}

func TestEngine_GetCurrentState(t *testing.T) {
	wfe := wf.NewEngine()

	s0 := wfe.GetInitialState()
	assert.NotNil(t, s0)
	assert.EqualValues(t, s0, wfe.GetCurrentState())

	s1, _ := wfe.RegisterState(S1)
	assert.NotNil(t, s1)

	err := wfe.RegisterEvent(s0, E1, s1)
	assert.NoError(t, err)

	nextState, err := wfe.ProcessEvent(E1)
	require.NoError(t, err)
	assert.EqualValues(t, nextState, wfe.GetCurrentState())
}

func TestEngine_GetState_When_Error(t *testing.T) {
	wfe := wf.NewEngine()
	s1 := wfe.GetState("FAKE STATE")
	assert.Nil(t, s1)
}

func TestEngine_ProcessEvent_When_Error(t *testing.T) {
	wfe := wf.NewEngine()

	s0 := wfe.GetInitialState()
	s1, _ := wfe.RegisterState(S1)

	err := wfe.RegisterEvent(s0, E1, s1)
	require.NoError(t, err)

	_, err = wfe.ProcessEvent(E2)
	require.Error(t, err)
	assert.EqualValues(t, "event \"EVENT_2\" is not defined for the current state \"INITIAL_STATE\"", err.Error())
}

func TestFullFlow(t *testing.T) {
	wfe := wf.NewEngine()

	s0 := wfe.GetInitialState()
	s1, _ := wfe.RegisterState(S1)
	s2, _ := wfe.RegisterState(S2)
	s3, _ := wfe.RegisterState(S3)
	s4, _ := wfe.RegisterState(S4)
	s5, _ := wfe.RegisterState(S5)

	err := wfe.RegisterEvent(s0, E1, s1)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s0, E2, s2)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s1, E3, s3)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s2, E3, s3)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s3, E4, s4)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s4, E5, s2)
	assert.NoError(t, err)

	err = wfe.RegisterEvent(s4, E6, s5)
	assert.NoError(t, err)

	nextState, err := wfe.ProcessEvent(E1)
	require.NoError(t, err)
	assert.EqualValues(t, s1, nextState)

	nextState, err = wfe.ProcessEvent(E3)
	require.NoError(t, err)
	assert.EqualValues(t, s3, nextState)

	nextState, err = wfe.ProcessEvent(E4)
	require.NoError(t, err)
	assert.EqualValues(t, s4, nextState)

	nextState, err = wfe.ProcessEvent(E6)
	require.NoError(t, err)
	assert.EqualValues(t, s5, nextState)

	assert.EqualValues(t, s5, wfe.GetCurrentState())
}
