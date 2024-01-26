package quotes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Quote struct {
	Author string `json:"author"`
	Source string `json:"source"`
	Text   string `json:"text"`
	Tags   string `json:"tags"`
}

type Tag struct {
	Key string `json:"key"`
	Val int    `json:"val"`
}

func (q *Quote) GetAuthor() string {
	if q.Author == "" {
		return "Неизвестный автор"
	}

	return q.Author
}

func (q *Quote) GetSource() string {

	return q.Source
}

func (q *Quote) GetText() string {
	return q.Text
}

func (q *Quote) GetTags() string {
	return q.Tags
}

func(q Quote) String() string {
    text := fmt.Sprintf("%s\n%s - %s\n",
		q.GetText(),
		q.GetAuthor(),
		q.GetSource(),
    )

    if len(q.GetTags()) > 0 {
        text += fmt.Sprintf("%s", q.GetTags())
    }

    return text
}

// Makes a generic resuquest to an api endpoint specified by argumennt
func NewRequestWithContext(r string)([]byte, error) {
    var data []byte

    API := os.Getenv("API")

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx,
        "GET",
        API + r,
        nil,
   )

   if err != nil {
       return data, fmt.Errorf("NewRequestWithContext failed: %w", err)
   }

   res, err := http.DefaultClient.Do(req)

   if err != nil {
       return data, fmt.Errorf("NewRequestWithContext: http.DefaultClient failed: %w", err)
   }

   defer res.Body.Close()

   data, err = io.ReadAll(res.Body)

   if err != nil {
       return data, fmt.Errorf("NewRequestWithContext: io.ReadAll failed: %w", err)
   }

   return data, nil
}

func GetRandomQuote()(*Quote, error) {
    var (
        data []byte
        err error
        quote = Quote{}
    )

    data, err = NewRequestWithContext("random")

    if err != nil {
        return &quote, fmt.Errorf("API: RandomQuote failed: %w", err)
    }

    err = json.Unmarshal(data, &quote)

    if err != nil {
        err = fmt.Errorf("API: RandomQuote fialed to unmarshal: %w", err)
    }

    return &quote, err
}

func GetTags()([]Tag, error) {

    var(
        data []byte
        err error
        tags []Tag
    )

    data, err = NewRequestWithContext("tags/10")
    if err != nil {
        return tags, err
    }


    err = json.Unmarshal(data, &tags)

    return tags, err
}

func GetQuoteByTag(t string)(*Quote, error) {
    var (
        data    []byte
        err     error
        quote   = Quote{}
    )

    //INFO: API doesn't return an error if case there are no quotes on this tag
    // should fix taht or work around it.

    t = strings.ToLower(t)
    t = strings.Trim(t, "\n\t ")

    data, err = NewRequestWithContext("tag/" + t)

    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(data, &quote)

    if quote.Text == "" {
        quote.Text = fmt.Sprintf("Похоже по тегу %s не найдено ни одной цитаты. Попробуйте другой.",
            t,
        )
        quote.Author = "администрация"
    }


    return &quote, err
}
