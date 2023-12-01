package wf_test

import (
	"testing"

	"github.com/MickStanciu/go-state/wf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_GetName(t *testing.T) {
	wfe, err := wf.NewEngine()
	require.NoError(t, err)
	require.NotNil(t, wfe)

	s0 := wfe.GetInitialState()
	assert.NotNil(t, s0)
	assert.EqualValues(t, "STATE_INITIAL", s0.GetName())
}
