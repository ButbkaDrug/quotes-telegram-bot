package bot

import (
	"fmt"
	"log"
	"net/http"

	"github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	"github.com/butbkadrug/advanced-telegram-bot-go/pkg/telegram"
	"github.com/butbkadrug/advanced-telegram-bot-go/pkg/witai"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//INFO: Bot package should be responsible for establishing connecion with
//telegram server and receiving / sending updates. All the processing should
// be done in the other packages.


type UpdateType int

const (
	TextMessageType UpdateType = iota
	VoiceMessageType
	AudioMessageType
	CallbackQueryType
	UnknownType
)

type bot struct {
	// go telegram bot api wrapper
	api *telegram.Telegram
	// wit.ai api wrapper
	wit *witai.Wit
}

type Update struct {
    tgbotapi.Update
}

func NewUpdate(u tgbotapi.Update) Update{
    return Update{
        u,
    }
}

func (u *Update) Type() UpdateType {

	if u.Message == nil && u.CallbackQuery != nil {
		return CallbackQueryType
	}

	if u.Message == nil {
		return UnknownType
	}

	if u.Message.Voice != nil {
		return VoiceMessageType
	}

	if u.Message.Audio != nil {
		return AudioMessageType
	}

	if u.Message.Text != "" {
		return TextMessageType
	}

	return UnknownType
}

// NewBotWithKey creates a new bot with a given key. If you don't have a key
// you can get it from @BotFather.
func NewBotWithKey(key string) (*bot, error) {
	var err error
	var bot = &bot{}

	ai := witai.NewWit()

	log.Println("Starting a bot...")

	api, err := telegram.NewTelegram(key)

	if err != nil {
		return bot, fmt.Errorf("NewBotWithKey: failed to start a bot: %w", err)
	}

	log.Printf("Bot %s is started\n", api.Info.UserName)

	bot.api = api
	bot.wit = ai

	return bot, nil
}

// Run starts listening for updates and handles them.
func (b *bot) Run() {

	updates := b.api.GetUpdatesChan()

	for update := range updates {
		b.UpdateHandler(NewUpdate(update))
	}
}

// UpdateHandler shoud be it's own package
//TODO: refactor UpdateHandler to be in separate packag

func (b *bot) UpdateHandler(update Update) {
    var id = update.FromChat().ID

	updateUserErr := userbase.UpdateLastVisited(id)

	if updateUserErr != nil {
		b.api.Logger.Error("failed to update user!", updateUserErr)
	}


    switch update.Type() {
    case TextMessageType:
        if update.Message.IsCommand(){
            b.MessageWithCommandUpdate(update)
            return
        }

        b.TextMessageUpdate(update)

    case VoiceMessageType:
        resp, err := b.VoiceMessageUpdate(update)

        if err != nil {
            b.api.SendInternalErrorMessage(id, err)
            return
        }

        b.api.Router(id, resp)



    default:
        b.api.SendBadRequestMessage(id)
    }

	// if title := user.Graduate(); title.Message != "" {
	//     err := b.api.SendTextMessage(update.FromChat().ID, title.Message)
	//     if err != nil {
	//         log.Println("Faild to send cheerful message: ", err)
	//     }
	// }
}

func (b *bot) MessageWithCommandUpdate(update Update) {
    var id = update.FromChat().ID

    switch update.Message.Command(){
    case "start", "help":
        b.api.SendWelcomeMessage(id)
    case "tags":
        b.api.SendTagsList(id)
    case "tag":
        //
    case "stats", "statistic":
        b.api.SendUserStats(id)
    case "random":
        b.api.SendRandomQuote(id)
    default:
        b.api.SendBadRequestMessage(id)
    }
}

func (b *bot) TextMessageUpdate(update Update) {

    var id = update.FromChat().ID
    resp, err := b.wit.HandleTextMessage(update.Message.Text)

    if err != nil {
        b.api.SendConnectionErrorMessage(id, err)
    }

    b.api.Router(id, resp)
}

//NOTE: I need user input proccesser with will respond with a command. Which
// then will be taken, validated and acted upon. May be bot will be an input
// proccesser package or maybe it should be it's own. And then command package
// will link user input with actual actions that bot needs to take.

func (b *bot) VoiceMessageUpdate(update Update) (telegram.Response, error){
	//TODO: refactor this to be in separate package
	//EXAMPLE: b.handler.VoiceMessageHandler(update)

	var wr witai.WitResponse

    link, err := b.api.GetFileLink(update.Message.Voice.FileID)

    if err != nil {
        return wr, err
    }


	resp, err := http.Get(link)

	if err != nil {
        return wr, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
        return wr, err
	}

    wr, err = b.wit.ParseAudio(resp.Body)

    if err != nil {
        return wr, err
    }

    return wr, err
}
