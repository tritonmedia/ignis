package telegram

import (
	"bytes"
	"fmt"
	"strings"

	router "github.com/tritonmedia/ignis/pkg/router"
)

func listStageFn(s *router.Scene) error {
	m, err := tr.ListMedia()
	if err != nil {
		return fmt.Errorf("failed to list media: %v", err)
	}

	var b bytes.Buffer

	fmt.Fprintln(&b, "🎬 "+locale.Strings.LISTHEADER)
	for i, media := range m {
		prefix := "├"
		if i+1 == len(m) {
			prefix = "└"
		}

		st := strings.Title(strings.ToLower(media.Status.String()))
		ty := strings.Title(strings.ToLower(media.Type.String()))
		if ty == "Tv" {
			// HACK
			ty = "TV"
		}
		fmt.Fprintf(&b, "%s *%s*  (Type: %s / Status: %s)\n", prefix, media.Name, ty, st)
	}

	if len(m) == 0 {
		fmt.Fprintln(&b, locale.Strings.LISTEMPTY)
	}

	s.Done()
	return s.Message.Send(b.String(), s.Message.GetID())
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
