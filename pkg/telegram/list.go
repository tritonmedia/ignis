package telegram

import (
	"fmt"

	router "github.com/tritonmedia/ignis/pkg/router"
)

func listStageFn(s *router.Scene) error {
	m, err := tr.ListMedia()
	if err != nil {
		return fmt.Errorf("failed to list media: %v", err)
	}

	resp := "🎬 " + locale.Strings.LISTHEADER + "\n"
	for i, media := range m {
		prefix := "├"
		if i+1 == len(m) {
			prefix = "└"
		}

		str := fmt.Sprintf("%s *%s* (Status %s)\n", prefix, media.Name, media.Status.String())
		resp += str
	}

	if len(m) == 0 {
		resp += locale.Strings.LISTEMPTY
	}

	s.Done()
	return s.Message.Send(resp, s.Message.GetID())
}

func listStage() *router.Stage {
	sc, _ := router.NewStage("listFn", listStageFn)
	return sc
}

func newListScene() *router.Scene {
	s := router.NewScene("list")
	s.Add(listStage())
	return s
}
