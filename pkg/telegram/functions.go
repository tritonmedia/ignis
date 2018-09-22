package telegram

import (
	"log"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/tritonmedia/ignis/pkg/state"
	"gopkg.in/telegram-bot-api.v4"
)

// stageInit is run when a user firsts contacts the bot, or doesn't have a current stage (context)
func stageInit(msg *tgbotapi.Message, u *state.User, c *cache.Cache) (string, error) {
	v, found := c.Get("lastRunTime")
	c.Set("lastRunTime", time.Now(), cache.NoExpiration)

	if found == false {
		log.Print("[functions/init] no last runtime :scream_cat:")
		return "Hello! Try talking to me later ....", nil
	}

	lastRunTime := v.(time.Time)
	return "You last talked to me at: " + lastRunTime.UTC().Format(time.RFC3339), nil
}

func register() {
	registerFunc(stageInit, "init")
}
