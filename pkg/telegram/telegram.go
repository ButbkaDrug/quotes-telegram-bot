package telegram

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/butbkadrug/advanced-telegram-bot-go/internal/userbase"
	"github.com/butbkadrug/advanced-telegram-bot-go/pkg/quotes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandType int

const (
	Start CommandType = iota
	RandomQuote
	QuoteWithTag
	TagsList
	Stats
	Unknown
)

type Telegram struct {
	Info   *tgbotapi.User
	api    *tgbotapi.BotAPI
	token  string
	Logger *slog.Logger
}

type Response interface {
	Command() CommandType
	Arguments() []string
}

// NewTelegram creates a new telegram bot with a given token.
func NewTelegram(token string) (*Telegram, error) {

	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, fmt.Errorf("NewTelegram: failed to start a bot: %w", err)
	}

	me, err := api.GetMe()

	if err != nil {
		return nil, fmt.Errorf("NewTelegram: failed to get a bot info: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Telegram{
		Info:   &me,
		api:    api,
		token:  token,
		Logger: logger,
	}, nil
}

// MESSAGE SENDING METHODS

func (t *Telegram) SendTextMessage(id int64, text string) error {
	message := tgbotapi.NewMessage(id, text)

	_, err := t.api.Send(message)

	return err
}

// Generic Error Message that can be used to send error messages to the user.
func (t *Telegram) SendErrorMessage(id int64, err error, args ...any) error {
	t.Logger.Error(err.Error(), args...)
	_, e := t.api.Send(tgbotapi.NewMessage(id, err.Error()))
	return e
}

// Inform the user that something went wrong on the server side.
func (t *Telegram) SendInternalErrorMessage(id int64, err error) error {
	errMsg := `Изивините, произошла внутренняя ошибка. Попробуйте позже.
Если проблема не исчезнет, свяжитесь с администратором.`

	return t.SendErrorMessage(id, errors.New(errMsg), err)
}

// Informs the user about bed request.
func (t *Telegram) SendBadRequestMessage(id int64) error {
	errMsg := `Я всего лишь скромынй плод относительно плохой инженерии.
Не уверен, что я смогу выполнить ваш запрос. Если у вас есть предложения,
как помочь мне работать лучше, расскажите об этом адинистратору!`

	return t.SendErrorMessage(id, errors.New(errMsg))
}

// User will receive nice and friendly error message. Error will be logged
func (t *Telegram) SendConnectionErrorMessage(id int64, err error) error {
	errMsg := `Похоже наш домовой эльф снова отлучился помогать Гарри, Рону и
Гермионе. Поэтому мы не сможем обработаь ваш запрос прямо сейчас,
пожалуйста попробуйте позже.`

	return t.SendErrorMessage(id, errors.New(errMsg), err)

}

// END OF MESSAGE SENDING METHODS

// GetUpdatesChan starts listening for updates and returns a channel.
func (t *Telegram) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return t.api.GetUpdatesChan(u)
}

// Returns a link to a file hosted on telegram's server. You can use it to
// download the file or pass it to another service
func (t *Telegram) GetFileLink(fileid string) (string, error) {

	file, err := t.api.GetFile(tgbotapi.FileConfig{
		FileID: fileid,
	})

	if err != nil {
		t.Logger.Error("file request error", err)
		return "", err
	}
	return file.Link(t.token), nil
}

func (t *Telegram) SendWelcomeMessage(id int64) {

	text := `Привет! 👋 Я твой персональный бот для цитат.
Я могу поделиться с тобой интересными цитатами. Вот список доступных комманд:

/start - начать работу с ботом
/random - получить случайную цитату
/tags - получить список тегов
/tag <tag> - получить случайную цитату с тегом <tag>
/stats - получить статистику
С любовью,
Твой Бот 🤖
`
	err := t.SendTextMessage(id, text)

	if err != nil {
		t.Logger.Error("failed to send welcome message", err)
	}

}

func (t *Telegram) SendRandomQuote(id int64) {
	var quote = &quotes.Quote{}
	quote, err := quotes.GetRandomQuote()

	if err != nil {
		t.SendConnectionErrorMessage(id, err)
		return
	}

	err = t.SendTextMessage(id, fmt.Sprint(quote))

	if err != nil {
		t.Logger.Error("failed to send random quote: %w", err)
	}
}

func (t *Telegram) SendTagsList(id int64) {
	tags, err := quotes.GetTags()

	if err != nil {
		t.SendConnectionErrorMessage(id, err)
		return
	}

	var text string

	for _, tag := range tags {
		text += tag.Key + "\n"
	}

	err = t.SendTextMessage(id, text)

	if err != nil {
		t.Logger.Error("failed to send tags", err)
	}
}

func (t *Telegram) SendQuoteWithTag(id int64, tag string) {

	q, err := quotes.GetQuoteByTag(tag)

	if err != nil && q.Text == "" {
		t.SendConnectionErrorMessage(id, err)
		return
	}

	err = t.SendTextMessage(id, fmt.Sprint(q))

	if err != nil {
		t.Logger.Error("failed to send quote with tag", err)
	}

}

func (t *Telegram) SendUserStats(id int64){
    u, err := userbase.GetUser(id)

    if err != nil {
        t.SendConnectionErrorMessage(id, err)
        return
    }

    err = t.SendTextMessage(id, u.GetStatsString())

	if err != nil {
		t.Logger.Error("failed to send quote with tag", err)
	}
}

func (t *Telegram) Router(id int64, r Response) {

	switch r.Command() {
	case Start:
		t.SendWelcomeMessage(id)
	case RandomQuote:
		t.SendRandomQuote(id)
	case TagsList:
        t.SendTagsList(id)
	case QuoteWithTag:
        t.SendQuoteWithTag(id, r.Arguments()[0])
    case Stats:
        t.SendUserStats(id)
	default:
		t.SendBadRequestMessage(id)
	}
}
