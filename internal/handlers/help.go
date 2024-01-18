package handlers

import (
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var helpHandler = &Handler{
	Use:         "help",
	Visible:     true,
	Scorable:    false,
	Description: `Покажет тебе это сообщение.`,
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

        upd.Message.Text = "/start"

		return Root.Execute(upd, nil)
	},
}

func init() {
	Root.AddHandler(helpHandler)
}
