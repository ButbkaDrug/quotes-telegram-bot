package handlers

import (
	"github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
    // Use will be compared against command that user input and if match executed
    Use string

    // If true, command will be shown in a start menu
    Visible bool

    // If true user score will be updated when command executed
    Scorable bool

    // Self explenatory :)
    Description string

    // Run will be executed when user call a comand.
    // And result will be sent as a request to a td server.
    Run func(tblib.Update, ...[]interface{})(tblib.Chattable, error)
}

type Handlers struct {
    handlers []*Handler
}

func(h *Handlers) AddHandler(handler *Handler){
    h.handlers = append(h.handlers, handler)
}

func (h *Handlers) Execute(u tblib.Update, user *userbase.User) (tblib.Chattable, error) {
    // Before processing a command I want to interect with user a little bit
    // For example if it is a firt quote for today. I want to cheer him up with
    // a message.

    // Also, ranking system will be triggered here. If user fulfilled requirement
    // for a award and does not posess it yet. Congrat him with a message and a
    // picture etc.

    // For now let test this stuff only if it is me who makes a request.
    var handler = &Handler{}
    cmd := u.Message.Command()

    for _, h := range h.handlers {

        if h.Use == cmd {
            handler = h
        }
    }


    if handler.Run == nil {
        return tblib.NewMessage(
            u.Message.From.ID,
            "Извини, но такой команды не существует...\n\n /help - подскажет тебе наиболее полезные команды.",
        ), nil
    }

    if handler.Scorable {
        user.UpdateTotalRequests(u.FromChat().ID)
    }

    return handler.Run(u)
}


var Root = Handlers{}


func init() {

}

