package main

import (
	"log"
	"os"
    "github.com/butbkadrug/advanced-telegram-bot-go/internal/bot"
	"github.com/joho/godotenv"
)



func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load the enviroment: ", err)
	}

	key := os.Getenv("TESTBOT_API_KEY")

    bot, err := bot.NewBotWithKey(key)

    if err != nil {
        log.Fatal("Failed to create a bot: ", err)
    }

    bot.Run()
}

