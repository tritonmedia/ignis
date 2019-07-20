package telegram

import (
	"log"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/analysis"
	"github.com/tritonmedia/ignis/pkg/state"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// canProceed determines if we have a positive go-ahead or negative go-ahead
func canProceed(msg string) bool {
	a := analysis.ProceedAnalysis(msg)

	return a
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

	if !canProceed(msg.Text) {
		return `OK! I've went ahead and created the request. Here is it's link: `, nil
	}

	return `Sorry, I don't support adding URLs yet. But I've created the request, here's the link: ` + `
	
Let me know if I can create another request sometime.
`, nil
}
