package main

import (
	"log"
	"os"
    "flag"
    "github.com/butbkadrug/advanced-telegram-bot-go/internal/bot"
	"github.com/joho/godotenv"
)

func main() {
    var key string

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load the enviroment: ", err)
	}

    flag.StringVar(&key, "key", os.Getenv("QUOTE_DRAGON_API_KEY"), "Telegram bot API key")
    flag.Parse()

    bot, err := bot.NewBotWithKey(key)

    if err != nil {
        log.Fatal("Failed to create a bot: ", err)
    }

    bot.Run()
}

