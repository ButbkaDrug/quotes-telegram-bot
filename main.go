package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/butbkadrug/advanced-telegram-bot-go/internal/handlers"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type bot struct {
	api *tblib.BotAPI
}


func main() {
    fmt.Println(Root)

	var err error
	var bot bot



	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load the enviroment: ", err)
	}

	key := os.Getenv("TESTBOT_API_KEY")

	b, err := tblib.NewBotAPI(key)

	bot.api = b

	if err != nil {
		log.Fatal("Failed to start a bot: ", err)
	}

	u := tblib.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	for update := range updates {
		var request tblib.Chattable

		fmt.Println(update)

		if update.Message != nil && update.Message.IsCommand() {
            request, err = Root.Execute(update.Message.Command(), update)

            if err != nil {
                log.Println("Failed to execute the command: ", err)
            }
		}


		if request == nil {
			continue
		}

		if _, err := bot.api.Send(request); err != nil {
			log.Println("Failed to send the message: ", err)
		}

	}
}

