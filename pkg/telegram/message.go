package telegram

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"text/template"

	"github.com/tritonmedia/ignis/pkg/state"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

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

	log.Printf("[telegram/message] Send(<stripped>, %s)", ID)

	_, err = t.Bot.Send(m)
	return err
}

// GetID returns the current chat ID
func (t *TelegramMessage) GetID() string {
	return strconv.Itoa(int(t.Message.Chat.ID))
}

// Text returns the text of the message
func (t *TelegramMessage) Text() string {
	return t.Message.Text
}
