package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

    . "github.com/butbkadrug/advanced-telegram-bot-go/internal/models"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tagHandler = &Handler{
    Use: "tag",
    Visible: false,
    Description: "You can select get quotes by tag",
    Run: func(upd tblib.Update, args ...[]interface{})(tblib.Chattable, error){

        var quote Quote
        var text = "Get that quote on a specific topic:"
        var message = tblib.NewMessage(upd.Message.From.ID, text)

        tag := upd.Message.CommandArguments()

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

        if quote.Author == "" {quote.Author = "Неизвестный автор"}

        text += "\n\n" + quote.Author

        if quote.Source != "" {

            text += " - " + quote.Source
        }

        message.Text = text

       return message, nil
    },
}


func init(){
    Root.AddHandler(tagHandler)
}
