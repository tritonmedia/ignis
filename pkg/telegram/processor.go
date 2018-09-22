package telegram

import (
	"errors"
	"log"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/state"
	"gopkg.in/telegram-bot-api.v4"
)

// fn is a function table function, should take the message and a user as the input
type fn func(*tgbotapi.Message, *state.User, *cache.Cache) (string, error)

var functionTable map[string]fn

// processMessage runs an action when a message is recieved
func processMessage(msg *tgbotapi.Message, s *state.State, u *state.User) (*tgbotapi.MessageConfig, error) {
	if fn, ok := functionTable[u.Stage]; ok {
		stage := state.NewStageStorage(u, u.Stage)

		resp, err := fn(msg, u, stage)
		if err != nil {
			return nil, err
		}

		m := tgbotapi.NewMessage(msg.Chat.ID, resp)
		m.ReplyToMessageID = msg.MessageID

		return &m, nil
	}

	return nil, errors.New("Invalid stage")
}

func registerFunc(f fn, stageName string) error {
	log.Printf("registering stage: %s", stageName)
	functionTable[stageName] = f

	return nil
}
