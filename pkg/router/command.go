package router

import (
	"fmt"
	"log"
)

// CommandRouter is a scene router that works by reading associating
// one words (i.e !new -> new or /start -> start) and triggering a scene by that name
type CommandRouter struct {
	// gernators are a map of scene generators used to create
	// scene structs for each user
	generators map[string]CommandRouterSceneGenerator

	// currentUserScenes is a map that tracks which scene a user is currently on
	currentUserScenes map[string]string

	// userScenes maps a userID to scenes that the user has. Format looks like
	// map[userId]map[sceneName]*Scene
	userScenes map[string]map[string]*Scene
}

// CommandRouterSceneGenerator is a CommandRouter acceptable scene generator
type CommandRouterSceneGenerator func() *Scene

// NewCommandRouter returns a new command router
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		generators:        make(map[string]CommandRouterSceneGenerator),
		currentUserScenes: make(map[string]string),
		userScenes:        make(map[string]map[string]*Scene),
	}
}

// AddScene registers a scene generator to this router
func (c *CommandRouter) AddScene(command string, s CommandRouterSceneGenerator) {
	if c.generators == nil {
		c.generators = make(map[string]CommandRouterSceneGenerator)
	}

	c.generators[command] = s
}

// Process determines which stage to run for a user's message
func (c *CommandRouter) Process(command string, m Message) error {
	id := m.GetID()
	if c.userScenes[id] == nil {
		c.userScenes[id] = make(map[string]*Scene)
	}

	// if we have a current scene, reset the state if it's been marked as finished
	if c.currentUserScenes[id] != "" {
		scene := c.currentUserScenes[id]
		if c.userScenes[id][scene].Finished {
			log.Printf("[router/command] Process(): resetting finish scene '%s'", scene)
			c.userScenes[id][scene].Reset()
			c.currentUserScenes[id] = ""
		} else {
			// set command to the last known command we ran, the user is supplying arbitrary data
			// at this point
			command = c.currentUserScenes[id]
		}
	}

	// if we have no scene for this user, and command is empty then it's a runtime error
	if c.currentUserScenes[id] == "" && command == "" {
		log.Println("[router/command] Process(): missing a command, and not current in a scene")
		return fmt.Errorf("missing command")
	}

	// if we don't have a current scene and have command, then start one
	if c.currentUserScenes[id] == "" && command != "" {
		if _, ok := c.generators[command]; !ok {
			log.Println("[router/command] Process(): failed to start a new stage, command not found")
			return fmt.Errorf("command not found: %s", command)
		}
		c.userScenes[id][command] = c.generators[command]()
		c.currentUserScenes[id] = command
	}

	// if we're this far then we should be able to progress onto the next stage of the scenes
	return c.userScenes[id][command].Next(m)
}
