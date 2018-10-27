package telegram

import (
	"log"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/analysis"
	"github.com/tritonmedia/ignis/pkg/state"
	"github.com/tritonmedia/ignis/pkg/trello"
	"gopkg.in/telegram-bot-api.v4"
)

// canProceed determines if we have a positive go-ahead or negative go-ahead
func canProceed(msg string) bool {
	a := analysis.ProceedAnalysis(msg)

	return a
}

// stageInit is run when a user firsts contacts the bot, or doesn't have a current stage (context)
func stageInit(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	s.SetStage(u.ID, "create-media")

	return `
Hello, [@{{.User.Username}}](tg://user?id={{.User.ID}})!

You can talk to me to create media request cards on the Triton Media platform.
The board is located at https://trello.com/b/vIGH0IiL/media-board.

To create a card, start by telling me the name of the media you'd like to request:

	_KonoSuba_
	`, nil
}

// createMedia will set the create media
func createMedia(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	s.SetStage(u.ID, "is-movie")

	c.Set("media:title", msg.Text, time.Hour)

	return `
You want to request media called "{{.Message.Text}}", correct?

Reply with 'cancel' to cancel, or 'yes' to continue
	`, nil
}

// createMediaConfirm processes the create media response
func createMediaConfirm(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	isMovie := false
	if canProceed(msg.Text) {
		isMovie = true
	}

	c.Set("media:isMovie", isMovie, time.Hour)

	s.SetStage(u.ID, "source-precheck")

	return `
OK! I'll create a request for this media for you. Do you have a link to a source for this?

If so, reply with 'yes', or 'no' if you don't.
	`, nil
}

// isMovie asks the user if we're a movie or not, and processes the createMedia answer
func isMovie(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	if !canProceed(msg.Text) {
		s.SetStage(u.ID, "is-movie")
		return "OK, let's try again. Please tell me the name of the show / movie you'd like to request", nil
	}

	s.SetStage(u.ID, "create-media-confirm")

	return `
Is this a movie?
	`, nil
}

// sourcePrecheck will term or ask for more info
func sourcePrecheck(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	s.SetStage(u.ID, "init")

	r, d := c.Get("media:title")
	title := r.(string)
	if !d || title == "" {
		s.SetStage(u.ID, "init")
		return `Failed to process.`, nil
	}

	r, _ = c.Get("media:isMovie")
	isMovie := r.(bool)

	log.Printf("[telegram/functions:sourcePrecheck] create: title=%s,source=%s,isMovie=%s", title, "", strconv.FormatBool(isMovie))

	link, err := trello.CreateCard(title, "", isMovie)
	if err != nil {
		log.Printf("[telegram/functions:sourcePrecheck] failed to create card: %s", err.Error())
		return `
I'm sorry, I failed to create the card. Please try again later!
		`, nil
	}

	if !canProceed(msg.Text) {
		return `OK! I've went ahead and created the request. Here is it's link: ` + link, nil
	}

	return `Sorry, I don't support adding URLs yet. But I've created the request, here's the link: ` + link + `
	
Let me know if I can create another request sometime.
`, nil
}

func register() {
	registerFunc(stageInit, "init")
	registerFunc(createMedia, "create-media")
	registerFunc(isMovie, "is-movie")
	registerFunc(createMediaConfirm, "create-media-confirm")
	registerFunc(sourcePrecheck, "source-precheck")
}
