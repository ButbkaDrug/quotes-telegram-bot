package bot

import (
    "fmt"
    "log"
    "github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	. "github.com/butbkadrug/advanced-telegram-bot-go/internal/handlers"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	api *tblib.BotAPI
}

func NewBotWithKey(key string) (*bot, error) {
	var err error
    var bot = &bot{}

    log.Println("Starting a bot...")

    api, err := tblib.NewBotAPI(key)

	if err != nil {
        return bot, fmt.Errorf("NewBotWithKey: failed to start a bot: %w", err)
	}

    me, err := api.GetMe()

    if err != nil {
        log.Println("Failed to get a bot info: ", err)
    }

    log.Printf("Bot %s is started\n", me.UserName)


    bot.api = api

    return bot, nil
}

func(b *bot) Run() {
    var err error
	u := tblib.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

    log.Println("Starting to listen for updates...")

	for update := range updates {
        log.Printf("%s %s(%d) %v\n",
            update.Message.From.FirstName,
            update.Message.From.LastName,
            update.Message.From.ID,
            update.FromChat(),
        )

        err = userbase.UpdateUser(update.FromChat().ID)

        if err != nil {
            log.Println(err)
        }

        n, err := userbase.UsersCount()

        if err != nil {
            log.Println(err)
        }

        log.Printf("Dataase updated: %d users in database", n)
		var request tblib.Chattable

		if update.Message != nil && update.Message.IsCommand() {
            request, err = Root.Execute(update.Message.Command(), update)

            if err != nil {
                log.Println("Failed to execute the command: ", err)
                request = tblib.NewMessage(update.Message.From.ID, "Упс! 😕 Произошла ошибка при попытке получить цитату. Пожалуйста, попробуйте еще раз позже или используйте другие команды. Если проблема сохраняется, обратитесь к администратору бота.")
            }
		}


		if request == nil {
			continue
		}

		if _, err := b.api.Send(request); err != nil {
			log.Println("Failed to send the message: ", err)
		}

	}
}
