package handlers

import (
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/api"
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/models"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var RandomQuoteHandler = &Handler{
	Use:         "random",
    Visible: true,
	Description: "This command with give you a random quote from our collection",
	Run: func(upd tblib.Update, args ...[]interface{}) (tblib.Chattable, error) {

		var quote = Quote{}
        quote, err := api.RandomQuote()


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
