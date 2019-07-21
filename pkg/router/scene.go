package router

import (
	"fmt"
	"log"
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

	// NextStage is the next stage that we should go onto, set by Enter()
	NextStage string

	// Bailed signifies if the last stage bailed on leave, this prevents it from running again
	Bailed bool

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
	s.NextStage = stage
}

// Next enters a scene, calling OnLeave on the existing stage
// and calling OnEnter on the new one, then calling the function
func (s *Scene) Next(m Message) error {
	// if next stage is empty, then assume we're going to the entrypoint
	stage := s.NextStage
	if s.NextStage == "" {
		stage = s.EntryPoint
	}

	log.Printf("[router/scene] Next(): moving into stage '%s' from '%s'", s.NextStage, s.CurrentStage)

	// update the message we're processing, this allows OnLeave()
	// to process the next message before the next stage gets it
	s.Message = m

	// if we were on a stage, then execute it's OnLeave() fn, and ensure it's OK
	if s.CurrentStage != "" && !s.Bailed {
		log.Printf("[router/scene] Next(): executing OnLeave for stage %s", s.CurrentStage)
		bail, err := s.Stages[s.CurrentStage].OnLeave(s)
		if err != nil {
			return err
		}

		// we've been told to bail, so don't update our current stage and assume that the leave
		// function modified the state and will be changed on the next message
		if bail {
			s.Bailed = true
			log.Printf("[router/scene] Next(): stage '%s' OnLeave() bailed (next stage was set to '%s')", s.CurrentStage, s.NextStage)
			return nil
		}
	}

	if st, ok := s.Stages[stage]; ok {
		// set the current stage, and reset bailed values
		s.CurrentStage = stage
		s.Bailed = false

		log.Printf("[router/scene] Next(): executing Func for stage %s", st.Name)
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
