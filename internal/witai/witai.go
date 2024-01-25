package witai

import (
	"fmt"
	"io"
	"log"
	"os"

	tblib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	witai "github.com/wit-ai/wit-go"
)

type Wit struct{
    *witai.Client
}

func NewWit() *Wit {
    client := witai.NewClient(os.Getenv("WITAI_API_KEY"))
    return &Wit{
        client,
    }
}


type entity struct{
    name string
    value string
    confidance float64
}

type WitResponse struct {
    Intent entity
    Tags []entity
}

func(w WitResponse)assertCommand() string{
    switch w.Intent.value{
    case "quote_request":
        if len(w.Arguments()) > 0 {

            return "tag"
        }

        return "random"
    case "stats_request":
        return "stats"
    }

    return "start"
}

func(w WitResponse)Command() string {
   return w.assertCommand()
}

func(w WitResponse)Arguments()[]string {
    var result []string

    for _, t := range w.Tags {
        if t.name == "tag" {
            result = append(result, t.value)
        }
    }
    return result
}

func NewEntity(name string, m map[string]interface{}) (*entity, error) {
    err := fmt.Errorf("NewIntentFromMap faild: value doesn't exists!")
    v, ok := m["value"]

    if ! ok {
        return nil, err
    }

    value := v.(string)
    v, ok = m["confidence"]

    if !ok {
        return nil, err
    }

    confidance := v.(float64)
    return &entity{
        name: name,
        value: value,
        confidance: confidance,
    }, nil
}

func(w *Wit)parseWitResponse(r *witai.MessageResponse) (WitResponse, error) {
    var result = WitResponse{}

    for k, v := range r.Entities {
        values := v.([]interface{})
        vmap := values[0].(map[string]interface{})
        entity, err := NewEntity(k, vmap)

        if err != nil {
            return result, err
        }

        switch k {
        case "intent":
            result.Intent = *entity
        case "tag":
            result.Tags = append(result.Tags, *entity)
        case "stats":
            result.Tags = append(result.Tags, *entity)
        case "quote":
            result.Tags = append(result.Tags, *entity)
        default:
            log.Println("WitParser: ",  k, v)
        }
    }

    return result, nil
}

func(w *Wit)parseText(s string) (WitResponse, error){
    var result = WitResponse{}

    resp, err := w.Parse(&witai.MessageRequest{
        Query: s,
    })

    if err != nil  {
        return result, err
    }

    return w.parseWitResponse(resp)
}

func(w *Wit)HandleTextMessage(u tblib.Update) (WitResponse, error) {
    var result = WitResponse{}

    resp, err := w.parseText(u.Message.Text)

    if err != nil {
        return result, err
    }

    return resp, nil
}

func(w *Wit)ParseAudio(f io.Reader)(WitResponse, error) {
    var result = WitResponse{}


    resp, err := w.Speech(&witai.MessageRequest{
        Speech: &witai.Speech{
            File: f,
            ContentType: "audio/opus",
        },
    })

    if err != nil {

        log.Println("ERROR SENDING SPECH: ", err)

        return result, err
    }

    return w.parseWitResponse(resp)
}
