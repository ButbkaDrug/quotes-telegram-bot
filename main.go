package main

import (
	"fmt"
	"log"
	"os"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const(
    start = "start"
    help = "help"
)

func main(){

    if err := godotenv.Load(); err != nil {
        log.Fatal("Failed to load the enviroment: ", err)
    }

    key := os.Getenv("TESTBOT_API_KEY")

    bot, err := botapi.NewBotAPI(key)

    if err != nil {
        log.Fatal("Failed to start a bot: ", err)
    }

    u := botapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)



    for update := range updates {
        fmt.Println(update)
        var msg botapi.MessageConfig

        if update.Message == nil && update.Message.IsCommand() {

            CommandHandler(update)
            continue

        }


        buttonOne := botapi.NewKeyboardButton("/help")
        buttonTwo := botapi.NewKeyboardButton("/troll")
        row := botapi.NewKeyboardButtonRow(buttonOne, buttonTwo)
        replyKeyboard := botapi.NewReplyKeyboard(row)

        replyKeyboard.InputFieldPlaceholder = "Select command:"
        replyKeyboard.ResizeKeyboard = true
        replyKeyboard.OneTimeKeyboard = true

        fmt.Println("Callback data: ", update.CallbackData())
        m := botapi.NewMessage(update.Message.From.ID, "Message")

        m.ReplyMarkup = replyKeyboard

        if _, err := bot.Send(msg); err != nil {
            log.Println("Failed to send the message: ", err)
        }
    }
}

func CommandHandler(u botapi.Update) {
    var msg botapi.MessageConfig
    switch cmd := u.Message.Command(); cmd {
    case start, help:
        text := `Hello, welcome to the chatbot.
I'm a bot that can help you with your daily tasks.
You can use the following commands:
/help - to get help
/troll - to get trolled
/play - to play a game
/stop - to stop the game`

        msg = botapi.NewMessage(u.Message.From.ID, text)

    default:
        fmt.Println("Unknown command", cmd)

        msg = botapi.NewMessage(u.Message.From.ID, "I don't know that command")
        ibtn1 := botapi.NewInlineKeyboardButtonData("Option one", "1")
        ibtn2 := botapi.NewInlineKeyboardButtonData("Option two", "2")

        irow := botapi.NewInlineKeyboardRow(ibtn1, ibtn2)
        ikeyboard := botapi.NewInlineKeyboardMarkup(irow)
        msg.ReplyMarkup = ikeyboard
    }
}
