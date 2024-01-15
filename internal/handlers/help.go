package handlers

import (
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var helpHandler = &Handler {
	Use:         "help",
	Visible:     true,
    Description: `Покажет тебе это сообщение.`,
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

		return Root.Execute("start", upd)
	},
}

func init() {
	Root.AddHandler(helpHandler)
}
