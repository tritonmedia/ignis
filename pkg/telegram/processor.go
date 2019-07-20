package telegram

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"text/template"

	cache "github.com/patrickmn/go-cache"
	router "github.com/tritonmedia/ignis/pkg/router"
	"github.com/tritonmedia/ignis/pkg/state"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// fn is a function table function, should take the message and a user as the input
type fn func(*tgbotapi.Message, *state.User, *cache.Cache, *state.State) (string, error)

var functionTable map[string]fn

var userScenes map[string]map[int]*router.Scene

// TelegramMessage is a router.Message compatible Telegram provider
type TelegramMessage struct {
	// Message is the underlying telegram message
	Message *tgbotapi.Message

	// Bot is the underlying Telegram Bot
	Bot *tgbotapi.BotAPI

	// User is the user we have from state for this message
	User *state.User
}

// NewTelegramMessage creates a TelegramMessage
func NewTelegramMessage(msg *tgbotapi.Message, u *state.User, bot *tgbotapi.BotAPI) *TelegramMessage {
	return &TelegramMessage{
		Message: msg,
		User:    u,
		Bot:     bot,
	}
}

// Get returns the message contents
func (t *TelegramMessage) Get() interface{} {
	return t
}

// Send sends a message
func (t *TelegramMessage) Send(text, ID string) error {
	// we treat all responses like templates
	var tpl bytes.Buffer
	tmp := template.New("resp")
	tmp.Parse(text)
	tmp.Execute(&tpl, map[string]interface{}{
		"User":    t.User,
		"Message": t.Message,
	})

	i, err := strconv.Atoi(ID)
	if err != nil {
		return fmt.Errorf("failed to convert id to an int: %v", err)
	}

	m := tgbotapi.NewMessage(int64(i), tpl.String())
	m.ParseMode = "Markdown"

	log.Printf("[emulate] Send(%s, %s)", tpl.String(), ID)

	_, err = t.Bot.Send(m)
	return err
}

// GetID returns the current chat ID
func (t *TelegramMessage) GetID() string {
	return strconv.Itoa(int(t.Message.Chat.ID))
}

// processMessage runs an action when a message is recieved
func processMessage(msg *tgbotapi.Message, u *state.User, bot *tgbotapi.BotAPI) (*tgbotapi.MessageConfig, error) {
	// TODO(jaredallard): scope this to allow stages to be set for each
	if _, ok := userScenes["new"]; !ok {
		userScenes["new"] = make(map[int]*router.Scene)
	}

	if _, ok := userScenes["new"][u.ID]; !ok {
		log.Printf("[telegram/processMessage] creating scene for user %d", u.ID)
		userScenes["new"][u.ID] = newNewStage()
	}

	if msg.Command() == "new" {
		_, err := userScenes["new"][u.ID].Enter("", NewTelegramMessage(msg, u, bot))
		return nil, err
	}

	return nil, fmt.Errorf("Unrecognized command '%s'", msg.Command())
}

func registerFunc(f fn, stageName string) error {
	log.Printf("[telegram/processor:register] registering stage: %s", stageName)
	functionTable[stageName] = f

	return nil
}
