package handlers

import (
	"fmt"
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var statsHandler = &Handler{
	Use:         "stats",
	Visible:     true,
	Scorable:    false,
	Description: "Узнай сколько цитат ты уже прочел. И какой твой текущий 👑титул в цитатном мире",
	Run: func(ChatID int64, r HandlerRequest) (tblib.Chattable, error) {
		var text string
		var m = tblib.NewMessage(ChatID, text)

		user, err := userbase.GetUser(ChatID)

		if err != nil {
			m.Text = "Что-то пошло не так.. Не смог получить ваши данные. Пожалуйста, попробуйте позже..."
			return m, err
		}

		text = fmt.Sprintf("🗓️ Ты с нами уже: %s\n\n", user.TimeJoinedString())
		text += fmt.Sprintf("📚 Цитат прочитано: %d\n\n", user.Total_requests)
		text += fmt.Sprintf("🏅 Текущий титул: %s\n\n", user.CurrentTitle().Title)
		text += fmt.Sprintf("⏳ Цитат до следущего титула(%s): %d",
			user.NextTitle().Title,
			user.NextTitle().Request-user.Total_requests,
		)

		m.Text = text

		return m, nil
	},
}

func init() {
	Root.AddHandler(statsHandler)
}
