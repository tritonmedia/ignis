package router

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// Scene is a structured message response system
type Scene struct {
	// Name of this scene
	Name string

	// CurrentStage is the current stage we're on
	CurrentStage string

	// Stages are the stages of this scene
	Stages map[string]*Stage

	// Message is the current message we're processing
	Message Message

	// EntryPoint is the initial stage, evaluated to the be the first provided to Add()
	EntryPoint string

	// StageData holds data for a stage
	StageData map[string]*cache.Cache
}

// NewScene creates a Telegram scene
func NewScene(name string) *Scene {
	return &Scene{
		Name: name,
	}
}

// Add a stage to this scene
func (s *Scene) Add(st *Stage) {
	// if this is nil, then initialize it
	if s.Stages == nil {
		s.Stages = make(map[string]*Stage)
		s.StageData = make(map[string]*cache.Cache)
		s.EntryPoint = st.Name
	}

	s.Stages[st.Name] = st
	s.StageData[st.Name] = cache.New(cache.NoExpiration, 10*time.Minute)
}

// GetCache returns a cache object for a stage
func (s *Scene) GetCache(stageName string) (*cache.Cache, error) {
	// use the current stage, if we have one, and one wasn't provided
	if s.CurrentStage != "" && stageName == "" {
		stageName = s.CurrentStage
	}

	if c, ok := s.StageData[stageName]; ok {
		return c, nil
	}

	// if we're here, then we haven't found a cache and since Add() creates
	// these for us we can assume that it doesn't exist/isn't a valid stage.
	return nil, fmt.Errorf("failed to find a cache for stage '%s', does it exist?", stageName)
}

// Enter enters a scene, calling OnLeave on the existing stage
// and calling OnEnter on the new one, then calling the function
func (s *Scene) Enter(stage string, m Message) (string, error) {
	// if stage is empty & no current stage, assume we're starting
	if stage == "" && s.CurrentStage == "" {
		stage = s.EntryPoint
	} else if stage == "" && s.CurrentStage != "" {
		// if we have no stage, but we're on a stage then it's a runtime error to
		// try to assume the entrypoint
		return "", fmt.Errorf("missing argument to Enter()")
	}

	// if we have a current stage, execute leave on it
	if s.CurrentStage != "" {
		s.Stages[s.CurrentStage].OnLeave(s)
	}

	// update the message we're processing
	s.Message = m

	if st, ok := s.Stages[stage]; ok {
		st.OnEnter(s)
		r, err := st.Func(s)
		return r, err
	}

	// stage doesn't exist
	return "", fmt.Errorf("stage '%s' not found", stage)
}
