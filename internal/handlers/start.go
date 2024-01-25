package handlers

import (
	"fmt"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startHandler = &Handler{
	Use:         "start",
	Visible:     false,
	Scorable:    false,
	Description: "Начни с этой команды чтобы узнать, на что я способен🦸",
	Run: func(ChatID int64, r HandlerRequest) (tblib.Chattable, error) {

		var text, cmds string
		text = `Привет! 👋 Я твой персональный бот для цитат. Я могу поделиться с тобой интересными цитатами. Вот список доступных комманд:
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

		m := tblib.NewMessage(ChatID, text)

		m.ReplyMarkup = keyboard

		return m, nil
	},
}

func init() {
	Root.AddHandler(startHandler)
}
