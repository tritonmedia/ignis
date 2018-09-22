package telegram

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tritonmedia/ignis/pkg/config"
	"github.com/tritonmedia/ignis/pkg/state"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// NewListener starts a new listener.
func NewListener(config *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

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

	functionTable = make(map[string]fn)
	register()

	s, err := state.NewClient(stateFile)

	users, err := s.ListUsers()
	if err != nil {
		log.Print("failed to list users")
	} else {
		log.Print("list of current (known) users:")
		fmt.Println("ID\t Username\t Stage")
		for _, user := range users {
			fmt.Println(user.ID, "\t", user.Username, "\t", user.Stage)
		}
	}

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
			_, err := s.CreateUser(id, username)
			if err != nil {
				log.Printf("[state] WARN: failed to create user: %s (err: %s)", username, err.Error())
			}
		} else if err != nil {
			log.Printf("[state] failed to search for user: %s (err: %s)", username, err.Error())
			continue
		} else {
			log.Printf("[state] found user: %s (uid: %d)", username, user.ID)
		}

		resp, err := processMessage(update.Message, s, user)
		if err != nil {
			log.Printf("[processor] ERR: failed to respond to %s (err: %s)", update.Message.Text, err.Error())
			continue
		}

		bot.Send(resp)
	}

	return nil
}
