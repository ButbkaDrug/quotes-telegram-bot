package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	. "github.com/butbkadrug/advanced-telegram-bot-go/internal/models"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tagHandler = &Handler{
	Use:         "tag",
	Visible:     false,
	Scorable:    true,
	Description: "С помощью этой команды, можно получить цитату на определенную тему",
	Run: func(ChatID int64, r HandlerRequest) (tblib.Chattable, error) {

		var quote Quote
        var tag string
		var text = "С помощью этой команды можно получить цитату на определенную тему\n:"
		var message = tblib.NewMessage(ChatID, text)

        args := r.Arguments()

        if len(args) < 1 {

            return message, fmt.Errorf("Tag is not found..")

        }

        tag = args[0]


		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/tag/"+tag, nil)

		if err != nil {
			log.Println("Faild to create request: ", err)
			return message, err
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println("Faild to make request: ", err)
			return message, err
		}

		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)

		if err != nil {
			log.Println("Failed to get the data: ", err)
			return message, err
		}

		err = json.Unmarshal(data, &quote)

		if err != nil {
			log.Println("Faild unmarshalling some data: ", err)
			return message, err
		}

		text = quote.Text

		if quote.Author == "" {
			quote.Author = "Неизвестный автор"
		}

		text += "\n\n" + quote.Author

		if quote.Source != "" {

			text += " - " + quote.Source
		}


        if quote.Tags != "" {
            text += "\n" + quote.Tags
        }

		message.Text = text

		return message, nil
	},
}

func init() {
	Root.AddHandler(tagHandler)
}
