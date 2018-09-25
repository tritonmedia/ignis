package telegram

import (
	"strings"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/state"
	"gopkg.in/telegram-bot-api.v4"
)

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
	s.SetStage(u.ID, "create-media-confirm")

	return `
You want to request media called "{{.Message.Text}}", correct?

Reply with 'cancel' to cancel, or 'yes' to continue
	`, nil
}

// createMediaConfirm processes the create media response
func createMediaConfirm(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	if strings.ToLower(msg.Text) == "cancel" {
		s.SetStage(u.ID, "create-media")
		return "OK, let's try again.", nil
	}

	s.SetStage(u.ID, "source-precheck")

	return `
OK! I'll create a request for this media for you. Do you have a link to a source for this?

If so, reply with 'yes', or 'no' if you don't.
	`, nil
}

// sourcePrecheck will term or ask for more info
func sourcePrecheck(msg *tgbotapi.Message, u *state.User, c *cache.Cache, s *state.State) (string, error) {
	s.SetStage(u.ID, "init")

	if strings.ToLower(msg.Text) != "yes" {
		return `OK! I've went ahead and created the request. Here is it's link: <link>`, nil
	}

	return `
Sorry, I don't support adding URLs yet. But I've created the request, here's the link: <link>

Let me know if I can create another request sometime.
`, nil
}

func register() {
	registerFunc(stageInit, "init")
	registerFunc(createMedia, "create-media")
	registerFunc(createMediaConfirm, "create-media-confirm")
	registerFunc(sourcePrecheck, "source-precheck")
}
