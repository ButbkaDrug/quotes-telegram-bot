package handlers

import (
    "fmt"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startHandler = &Handler{
	Use:         "start",
	Visible:     false,
	Description: "Welcome to the bot. Enter the command and we'll rock",
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

		var text, cmds string
        text =`Привет! 👋 Я твой персональный бот для цитат. Я могу поделиться с тобой интересными цитатами. Вот список доступных комманд:
%s
С любовью,
Твой Бот 🤖
`

		keyboard := tblib.NewReplyKeyboard()

		for _, h := range Root.handlers {
			if !h.Visible {
				continue
			}
			button := tblib.NewKeyboardButton("/" + h.Use)
			row := tblib.NewKeyboardButtonRow(button)
			keyboard.Keyboard = append(keyboard.Keyboard, row)

			cmds += fmt.Sprintf("\n/%s - %s\n", h.Use, h.Description)
		}

        text = fmt.Sprintf(text, cmds)

		m := tblib.NewMessage(upd.Message.From.ID, text)

		m.ReplyMarkup = keyboard

		return m, nil
	},
}

func init() {
	Root.AddHandler(startHandler)
}
