package telegram

import (
	"fmt"

	router "github.com/tritonmedia/ignis/pkg/router"
)

func startStageFn(s *router.Scene) error {
	err := s.Message.Send(
		`
Hello, [@{{.User.Username}}](tg://user?id={{.User.ID}})!

You can talk to me to interface with the Triton Media platform (https://github.com/tritonmedia/triton).

Here are the various commands I support:
	
	/new - create a new media
	/list - list all media that is currently on the server
	`,
		s.Message.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	s.Done()
	return nil
}

func startStage() *router.Stage {
	sc, _ := router.NewStage("startFn", startStageFn)
	return sc
}

func newStartScene() *router.Scene {
	s := router.NewScene("start")
	s.Add(startStage())
	return s
}
