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

var tagsHandler = &Handler{
    Use: "tags",
    Visible: true,
    Description: "Select one of our popular tags:",
    Run: func(upd tblib.Update, args ...[]interface{})(tblib.Chattable, error){

        var tags Tags
        var text = "You can select one of your popular tags."
        var message = tblib.NewMessage(upd.Message.From.ID, text)

        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/tags/10", nil)

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

        err = json.Unmarshal(data, &tags)

        if err != nil {
            log.Println("Faild unmarshalling some data: ", err)
            return message, err
        }


        keyboard := tblib.NewOneTimeReplyKeyboard()

        i:= len(tags)-1

        for i >= 0 {
            tag := tags[i]
            button := tblib.NewKeyboardButton("/tag " + tag.Key)
            row := tblib.NewKeyboardButtonRow(button)
            keyboard.Keyboard = append(keyboard.Keyboard, row)
            i--
        }

        message.Text = text
        message.ReplyMarkup = keyboard

       return message, nil

    },
}


func init(){
    Root.AddHandler(tagsHandler)
}
