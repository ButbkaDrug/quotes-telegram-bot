package handlers

import (
	"fmt"

	"github.com/butbkadrug/advanced-telegram-bot-go/internal/api"
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/models"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var RandomQuoteHandler = &Handler{
	Use:         "random",
	Visible:     true,
	Description: `Получи вдохновляющую цитату! 📖 Бот выберет случайную цитату для тебя. Просто отправь эту команду, чтобы добавить каплю мудрости в свой день.`,
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

		var quote = &models.Quote{}
		quote, err := api.RandomQuote()

		if err != nil {
			return nil, fmt.Errorf("RandomQuoteHandler: failed to get a random quote: %w", err)
		}

		text := quote.Text

		if quote.Author == "" {
			quote.Author = "Неизвестный автор"
		}

		text += "\n\n" + quote.Author

		if quote.Source != "" {

			text += " - " + quote.Source
		}

		return tblib.NewMessage(upd.Message.From.ID, text), nil
	},
}

func init() {
	Root.AddHandler(RandomQuoteHandler)
}
