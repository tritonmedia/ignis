package telegram

import router "github.com/tritonmedia/ignis/pkg/router"

func confirmMediaFn(s *router.Scene) (string, error) {
	return "", s.Message.Send(
		"What media would you like to request?",
		s.Message.GetID(),
	)
}

func newConfirmMedia() *router.Stage {
	sc, _ := router.NewStage("confirmMedia", confirmMediaFn)
	return sc
}

func newNewStage() *router.Scene {
	s := router.NewScene("new")
	s.Add(newConfirmMedia())
	return s
}
