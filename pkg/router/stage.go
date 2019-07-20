package router

import (
	"fmt"
)

// Stage is a "step" in a scene
type Stage struct {
	// Name of this stage
	Name string

	// Func is the main function of a stage
	Func StageFunc

	// OnLeave is a function executed
	OnLeave StageTransitionFunc

	OnEnter StageTransitionFunc
}

// StageFunc is a function used for a stage
type StageFunc func(*Scene) error

// StageTransitionFunc is a func used for Leave/Enter events.
type StageTransitionFunc func(*Scene) error

// NewStage returns a stage instance that is able to be used in a stage
func NewStage(name string, f StageFunc) (*Stage, error) {
	// empty handler
	generic := func(s *Scene) error {
		return nil
	}

	if name == "" {
		return nil, fmt.Errorf("missing stage name")
	}

	return &Stage{
		Name:    name,
		Func:    f,
		OnLeave: generic,
		OnEnter: generic,
	}, nil
}

// SetFunc overrides the initial stage function
func (st *Stage) SetFunc(f StageFunc) {
	st.Func = f
}
