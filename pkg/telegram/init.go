package telegram

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tritonmedia/ignis/pkg/config"
	router "github.com/tritonmedia/ignis/pkg/router"
	"github.com/tritonmedia/ignis/pkg/state"
	triton "github.com/tritonmedia/tritonmedia.go"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var tr *triton.Client
var locale *LocalizationFile

// NewListener starts a new listener.
func NewListener(config *config.Config, localeName string) error {
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return err
	}

	log.Printf("[telegram/init] Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	stateFile := filepath.Join(wd, "ignis.db")

	s, err := state.NewClient(stateFile)

	users, err := s.ListUsers()
	if err != nil {
		log.Print("failed to list users")
	} else {
		for _, user := range users {
			s.SetStage(user.ID, "init")
			log.Printf("[telegram/init] reset stage: oldStage=%s,stage=init,username=%s,uid=%d", user.Stage, user.Username, user.ID)
		}
	}

	tr, err = triton.NewClient(config.Triton.Host, config.Triton.Token)
	if err != nil {
		return fmt.Errorf("failed to create triton client: %v", err)
	}
	log.Printf("[triton/init] client created")

	locale, err = LoadLocale(localeName)
	if err != nil {
		return fmt.Errorf("failed to load locale '%s': %v", localeName, err)
	}

	log.Printf("[localization/init] loaded locale '%s'", localeName)

	// create a command router
	c := router.NewCommandRouter()
	c.AddScene("new", newNewScene)
	c.AddScene("start", newStartScene)
	c.AddScene("list", newListScene)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		username := update.Message.From.UserName
		if username == "" {
			log.Printf("[state] no username, using calculated firstname + lastname")
			username = update.Message.From.FirstName + " " + update.Message.From.LastName
		}

		id := update.Message.From.ID

		log.Printf("[state] attempting to find user: %s (uid: %d)", username, id)
		user, err := s.GetUserByID(id)
		if s.IsNotFound(err) {
			log.Printf("[state] creating user: %s (uid: %d)", username, id)
			user, err = s.CreateUser(id, username)
			if err != nil {
				log.Printf("[state] WARN: failed to create user: %s (err: %s)", username, err.Error())
			}
		} else if err != nil {
			log.Printf("[state] failed to search for user: %s (err: %s)", username, err.Error())
			continue
		} else {
			log.Printf("[state] found user: %s (uid: %d)", username, user.ID)
		}

		log.Printf("[processor] DEBU: message: '%s'", update.Message.Text)

		err = c.Process(update.Message.Command(), NewTelegramMessage(update.Message, user, bot))
		if err != nil {
			log.Printf("[processor] ERR: failed to respond to %s (err: %s)", update.Message.Text, err.Error())

			m := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Message processing failed: %v", err))
			m.ReplyToMessageID = update.Message.MessageID

			bot.Send(m)
			continue
		}
	}

	return nil
}
