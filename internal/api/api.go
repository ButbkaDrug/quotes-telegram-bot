package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/butbkadrug/advanced-telegram-bot-go/internal/models"
)

type Api struct {}

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
       return data, err
   }

   res, err := http.DefaultClient.Do(req)

   if err != nil {
       return data, err
   }

   defer res.Body.Close()

   data, err = io.ReadAll(res.Body)

   if err != nil {
       return data, err
   }

   return data, nil
}

func RandomQuote()(*models.Quote, error) {
    var (
        data []byte
        err error
        quote = models.Quote{}
    )

    data, err = NewRequestWithContext("random")

    if err != nil {
        return &quote, err
    }

    err = json.Unmarshal(data, &quote)

    return &quote, err
}

func Tags()(*models.Tags, error) {

    var(
        data []byte
        err error
        tags models.Tags
    )

    data, err = NewRequestWithContext("tags/10")
    if err != nil {
        return &tags, err
    }

    err = json.Unmarshal(data, &tags)

    return &tags, err
}

func Tag(t string)(*models.Quote, error) {
    var (
        data    []byte
        err     error
        quote   *models.Quote
    )

    data, err = NewRequestWithContext("tag/" + t)

    if err != nil {
        return quote, err
    }

    err = json.Unmarshal(data, quote)

    return quote, err
}
