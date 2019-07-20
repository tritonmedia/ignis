package telegram

import (
	"fmt"
	"log"
	"strings"

	router "github.com/tritonmedia/ignis/pkg/router"
)

// ---
// Stage Functions
// ---

func getMediaStageFn(s *router.Scene) error {
	err := s.Message.Send(
		"What media would you like to request?",
		s.Message.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("isMovie")
	return nil
}

func getMediaStageLeaveFn(s *router.Scene) error {
	s.Cache.SetDefault("name", s.Message.Text())

	log.Printf("[scene/new]: getMediaLeave(): media name is '%s'", s.Message.Text())
	return nil
}

func isMovieStageFn(s *router.Scene) error {
	err := s.Message.Send(
		"Is this media a movie?",
		s.Message.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("confirmMeda")
	return nil
}

func isMovieStageLeaveFn(s *router.Scene) error {
	resp := strings.ReplaceAll(" ", strings.ToLower(s.Message.Text()), "")
	if resp == "yes" {
		// TODO(jaredallard): map to proto here
		s.Cache.SetDefault("type", "movie")
	} else {
		s.Cache.SetDefault("type", "tv")
	}

	return nil
}

func confirmMediaStageFn(s *router.Scene) error {
	v, _ := s.Cache.Get("name")
	name := v.(string)
	v, _ = s.Cache.Get("type")
	mediaType := v.(string)

	err := s.Message.Send(
		fmt.Sprintf("Going to create a request for a the media '%s' type '%s'.\nContinue?", name, mediaType),
		s.Message.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("uploadMedia")
	return nil
}

func confirmMedaStageLeaveFn(s *router.Scene) error {
	resp := strings.ReplaceAll(" ", strings.ToLower(s.Message.Text()), "")
	if resp != "yes" {
		return s.Message.Send(
			"Aborting.",
			s.Message.GetID(),
		)
	}

	return nil
}

func uploadMediaStageFn(s *router.Scene) error {
	s.Done()
	return nil
}

// ----
// Stages
// ----

func getMediaStage() *router.Stage {
	sc, _ := router.NewStage("getMedia", getMediaStageFn)
	sc.OnLeave = getMediaStageLeaveFn
	return sc
}

func isMovieStage() *router.Stage {
	sc, _ := router.NewStage("isMovie", isMovieStageFn)
	sc.OnLeave = isMovieStageLeaveFn
	return sc
}

func confirmMediaStage() *router.Stage {
	sc, _ := router.NewStage("confirmMeda", confirmMediaStageFn)
	sc.OnLeave = confirmMedaStageLeaveFn
	return sc
}

func uploadMediaStage() *router.Stage {
	sc, _ := router.NewStage("uploadMedia", uploadMediaStageFn)
	return sc
}

func newNewScene() *router.Scene {
	s := router.NewScene("new")
	s.Add(getMediaStage())
	s.Add(isMovieStage())
	s.Add(confirmMediaStage())
	s.Add(uploadMediaStage())
	return s
}
