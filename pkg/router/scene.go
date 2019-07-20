package router

import (
	"fmt"
	"runtime"
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

	// Cache is the cache for this scene, used for storing data between stages
	Cache *cache.Cache

	// Finished, is this scene finished?
	Finished bool
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
		s.EntryPoint = st.Name
	}

	s.Stages[st.Name] = st
	s.Cache = cache.New(cache.NoExpiration, 10*time.Minute)
}

// Reset resets a stage to the initial entrypoint
func (s *Scene) Reset() {
	s.CurrentStage = s.EntryPoint
	s.Cache = cache.New(cache.NoExpiration, 10*time.Minute)

	// TODO(jaredallard): investigate impact of running GC
	// trigger GC to cleanup old caches
	runtime.GC()
}

// Enter is meant to be used by a stage function to
// TODO(jaredallard): decide if we want to return an error here if stage doesn't exist
func (s *Scene) Enter(stage string) {
	s.CurrentStage = stage
}

// Next enters a scene, calling OnLeave on the existing stage
// and calling OnEnter on the new one, then calling the function
func (s *Scene) Next(m Message) error {
	// if current stage is empty, then assume we're going to the entrypoint
	stage := s.CurrentStage
	if s.CurrentStage == "" {
		stage = s.EntryPoint
	}

	// update the message we're processing, this allows OnLeave()
	// to process the next message before the next stage gets it
	s.Message = m

	// if we have a current stage, execute leave on it
	if s.CurrentStage != "" {
		s.Stages[s.CurrentStage].OnLeave(s)
	}

	if st, ok := s.Stages[stage]; ok {
		st.OnEnter(s)
		return st.Func(s)
	}

	// stage doesn't exist
	return fmt.Errorf("stage '%s' not found", stage)
}

// Done marks a scene as finished
func (s *Scene) Done() {
	s.Finished = true
}
