package models

type Quote struct {
    Author  string  `json:"author"`
    Source  string  `json:"source"`
    Text    string  `json:"text"`
    Tags    string  `json:"tags"`
}

type Tag struct {
    Key string  `json:"key"`
    Val int     `json:"val"`
}

type Tags []Tag

