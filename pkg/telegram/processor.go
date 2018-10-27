package telegram

import (
	"bytes"
	"errors"
	"log"
	"text/template"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/state"
	"gopkg.in/telegram-bot-api.v4"
)

// fn is a function table function, should take the message and a user as the input
type fn func(*tgbotapi.Message, *state.User, *cache.Cache, *state.State) (string, error)

var functionTable map[string]fn

// processMessage runs an action when a message is recieved
func processMessage(msg *tgbotapi.Message, s *state.State, u *state.User) (*tgbotapi.MessageConfig, error) {
	if fn, ok := functionTable[u.Stage]; ok {
		stage := state.NewStageStorage(u)

		resp, err := fn(msg, u, stage, s)
		if err != nil {
			return nil, err
		}

		var tpl bytes.Buffer
		t := template.New("resp")
		t.Parse(resp)
		t.Execute(&tpl, map[string]interface{}{
			"User":    u,
			"Message": msg,
		})

		m := tgbotapi.NewMessage(msg.Chat.ID, tpl.String())
		m.ParseMode = "Markdown"

		return &m, nil
	}

	return nil, errors.New("Invalid stage")
}

func registerFunc(f fn, stageName string) error {
	log.Printf("[telegram/processor:register] registering stage: %s", stageName)
	functionTable[stageName] = f

	return nil
}
