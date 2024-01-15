package handlers

import (
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startHandler = &Handler{
	Use:         "start",
	Visible:     false,
	Description: "Welcome to the bot. Enter the command and we'll rock",
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

		var text string
		text = "Hello, welcome to the quotes. Here is some commands that you can explore:\n\n"

		keyboard := tblib.NewReplyKeyboard()

		for _, h := range Root.handlers {
			if !h.Visible {
				continue
			}
			button := tblib.NewKeyboardButton("/" + h.Use)
			row := tblib.NewKeyboardButtonRow(button)
			keyboard.Keyboard = append(keyboard.Keyboard, row)

			text += h.Use + " - " + h.Description + "\n"
		}

		m := tblib.NewMessage(upd.Message.From.ID, text)

		m.ReplyMarkup = keyboard

		return m, nil
	},
}

func init() {
	Root.AddHandler(startHandler)
}
