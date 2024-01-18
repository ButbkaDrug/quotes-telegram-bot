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
	u := tblib.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

    log.Println("Starting to listen for updates...")

	for update := range updates {
        b.UpdateHandler(update)
	}
}

func(b *bot) UpdateHandler(update tblib.Update) {
    var err error

    user, err := userbase.GetUser(update.FromChat().ID)

    if err != nil {
        log.Println("Can't fetch user: ", err)
        return
    }

    err = user.UpdateLastVisited(update.FromChat().ID)

    if err != nil {
        log.Println(err)
        return
    }

    var request tblib.Chattable

    if update.Message != nil && update.Message.IsCommand() {
        request, err = Root.Execute(update, user)

        if err != nil {
            log.Println("Failed to execute the command: ", err)
            request = tblib.NewMessage(update.Message.From.ID, "Упс! 😕 Произошла ошибка при попытке получить цитату. Пожалуйста, попробуйте еще раз позже или используйте другие команды. Если проблема сохраняется, обратитесь к администратору бота.")
        }
    }


    if request == nil {
        return
    }

    if _, err := b.api.Send(request); err != nil {
        log.Println("Failed to send the message: ", err)
    }

    if title := user.Graduate(); title.Message != "" {
        m := tblib.NewMessage(
            update.FromChat().ID,
            title.Message,
        )

        _, err := b.api.Send(m)
        if err != nil {
            log.Println("Faild to send cheerful message: ", err)
        }
    }
}
