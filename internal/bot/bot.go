package bot

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	. "github.com/butbkadrug/advanced-telegram-bot-go/internal/handlers"
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/witai"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	api *tblib.BotAPI
    wit *witai.Wit
}
type Bupdate struct {
    tblib.Message
}

func(u Bupdate) Command()string {
    return u.Message.Command()
}

func(u Bupdate) Arguments()[]string{
    return strings.Split(u.Message.CommandArguments(), " ")
}

func NewBotWithKey(key string) (*bot, error) {
	var err error
    var bot = &bot{}

    ai := witai.NewWit()

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
    bot.wit = ai

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
    log.Printf("New update: %+v\n", update)

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

    if update.Message.Voice != nil {
        var wr witai.WitResponse
        file, err := b.api.GetFile(tblib.FileConfig{
            FileID: update.Message.Voice.FileID,
        })

        if err != nil {
            log.Println("Can't request file: ", err)
        }

        link := file.Link(b.api.Token)

        resp, err := http.Get(link)

        if err != nil {
            log.Println("Can't get file: ", err)
        }

        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK{
            wr, err = b.wit.ParseAudio(resp.Body)

            if err != nil {
                log.Println(err)
                return
            }
        }
        r, err := Root.Execute(update.FromChat().ID, wr, user)

        if err != nil {
            log.Println("Faild fuilfil free text: ", err)
        } else {
            request = r
        }
    }

    if update.Message != nil && update.Message.IsCommand() {
        u := Bupdate{*update.Message}
        request, err = Root.Execute(
            update.FromChat().ID,
            u,
            user,
        )

        if err != nil {
            log.Println("Failed to execute the command: ", err)
            request = tblib.NewMessage(update.Message.From.ID, "Упс! 😕 Произошла ошибка при попытке получить цитату. Пожалуйста, попробуйте еще раз позже или используйте другие команды. Если проблема сохраняется, обратитесь к администратору бота.")
        }
    }

    if update.Message != nil && !update.Message.IsCommand() && update.Message.Text != ""{

        c, err := b.wit.HandleTextMessage(update)

        if err != nil {
            log.Println("Faild to process text(UpdateHandler):", err)
        }

        r, err := Root.Execute(update.FromChat().ID, c, user)

        if err != nil {
            log.Println("Faild fuilfil free text: ", err)
        } else {
            request = r
        }
    }



    if request == nil {
        request = tblib.NewMessage(update.Message.From.ID, "Упс! 😕 Произошла ошибка при попытке получить цитату. Пожалуйста, попробуйте еще раз позже или используйте другие команды. Если проблема сохраняется, обратитесь к администратору бота.")
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
