package telegram

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	router "github.com/tritonmedia/ignis/pkg/router"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// ---
// Stage Functions
// ---

func getMediaStageFn(s *router.Scene) error {
	err := s.Message.Send(locale.Strings.GETMEDIAGETNAME, s.Message.GetID())
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("isMovie")
	return nil
}

func getMediaStageLeaveFn(s *router.Scene) (bool, error) {
	s.Cache.SetDefault("name", s.Message.Text())

	log.Printf("[scene/new]: getMediaLeave(): media name is '%s'", s.Message.Text())

	m, err := tr.ListMedia()
	if err == nil {
		log.Printf("[scene/new] getMediaLeave(): checking %d media for potential duplicates", len(m))
		matches := make([]string, 0)
		for _, media := range m {
			// log.Printf("[scene/new] getMediaLeave()::duplicates: '%s' ~~ '%s'", media.Name, s.Message.Text())
			if fuzzy.Match(s.Message.Text(), media.Name) {
				matches = append(matches, media.Name)
			}
		}

		log.Printf("[scene/new] getMediaLeave(): found %d fuzzy matches", len(matches))

		if len(matches) != 0 {
			resp := locale.Strings.DUPLICATESHEADER + "\n\n"
			for i, match := range matches {
				resp = resp + fmt.Sprintf(" *%d*. %s\n", i+1, match)
			}
			resp += "\n" + locale.Strings.DUPLICATESFOOTER

			err := s.Message.Send(resp, s.Message.GetID())
			if err != nil {
				log.Printf("[scene/new] getMediaLeave(): failed to send response: %v", err)
				return false, fmt.Errorf("failed to send response: %v", err)
			}

			s.Enter("duplicate")

			// bail because we just changed the next scene in a leave function
			return true, nil
		}
	} else {
		log.Printf("[scene/new] WARN getMediaLeave(): failed to list media: %v", err)
	}

	return false, nil
}

func duplicateStageFn(s *router.Scene) error {
	if Proceed(s.Message.Text()) {
		s.Message.Send(locale.Strings.DUPLICATESCANCEL, s.Message.GetID())
		s.Done()
		return nil
	}

	s.Enter("isMovie")
	return s.Next(s.Message)
}

func isMovieStageFn(s *router.Scene) error {
	err := s.Message.Send(
		locale.Strings.ISMOVIEGETTYPE,
		s.Message.GetID(),
	)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("confirmMedia")
	return nil
}

func isMovieStageLeaveFn(s *router.Scene) (bool, error) {
	if Proceed(s.Message.Text()) {
		// TODO(jaredallard): map to proto here
		s.Cache.SetDefault("type", "movie")
	} else {
		s.Cache.SetDefault("type", "tv")
	}

	return false, nil
}

func confirmMediaStageFn(s *router.Scene) error {
	v, _ := s.Cache.Get("name")
	name := v.(string)
	v, _ = s.Cache.Get("type")
	mediaType := v.(string)

	var tpl bytes.Buffer
	tmp := template.New("inline")
	tmp.Parse(locale.Strings.CONFIRMMEDIAASK)
	tmp.Execute(&tpl, map[string]interface{}{
		"name": name,
		"type": mediaType,
	})

	err := s.Message.Send(tpl.String(), s.Message.GetID())
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	s.Enter("uploadMedia")
	return nil
}

func confirmMedaStageLeaveFn(s *router.Scene) (bool, error) {
	if !Proceed(s.Message.Text()) {
		return false, s.Message.Send(locale.Strings.GENERALABORT, s.Message.GetID())
	}

	err := s.Message.Send(locale.Strings.CONFIRMMEDIALEAVE, s.Message.GetID())
	if err != nil {
		return false, fmt.Errorf("failed to send response: %v", err)
	}

	// TODO(jaredallard): create actual request here
	return false, nil
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

func duplicateStage() *router.Stage {
	sc, _ := router.NewStage("duplicate", duplicateStageFn)
	return sc
}

func confirmMediaStage() *router.Stage {
	sc, _ := router.NewStage("confirmMedia", confirmMediaStageFn)
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
	s.Add(duplicateStage())
	s.Add(isMovieStage())
	s.Add(confirmMediaStage())
	s.Add(uploadMediaStage())
	return s
}
