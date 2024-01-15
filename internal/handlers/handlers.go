package handlers

import (
	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
    // Use will be compared against command that user input and if match executed
    Use string

    // If true, command will be shown in a start menu
    Visible bool

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

func (h *Handlers) Execute(cmd string, u tblib.Update) (tblib.Chattable, error) {

    var handler = &Handler{}

    for _, h := range h.handlers {

        if h.Use == cmd {
            handler = h
        }
    }


    if handler.Run == nil {
        return tblib.NewMessage(u.Message.From.ID, "Sorry, command does't exist"), nil
    }

    return handler.Run(u)
}


var Root = Handlers{}


func init() {

}

