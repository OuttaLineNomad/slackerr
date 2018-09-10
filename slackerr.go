package slackerr

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Error customizes erros before the are sent from package.
type Error struct {
	Func string
	Err  error
}

// Fields sub structs for SendMsg.
type Fields struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Actions sub structs for SendMsg.
type Actions struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Style string `json:"style"`
	URL   string `json:"url"`
}

// Attachments sub structs for SendMsg.
type Attachments struct {
	Fallback   string    `json:"fallback"`
	Title      string    `json:"title"`
	AuthorName string    `json:"author_name"`
	AuthorIcon string    `json:"author_icon"`
	Color      string    `json:"color"`
	Fields     []Fields  `json:"fields"`
	Actions    []Actions `json:"actions"`
}

// SendMsg holds the payload for slack errors.
type SendMsg struct {
	Text        string        `json:"text"`
	Attachments []Attachments `json:"attachments"`
}

func (er *Error) Error() string {
	return `slackerr: ` + er.Func + `: ` + er.Err.Error()
}

// Send sends messages to slak.
func Send(channel string, pld *SendMsg, custom *http.Client) error {

	b, err := json.Marshal(pld)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", channel, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	clint := http.DefaultClient
	if custom != nil {
		clint = custom
	}
	res, err := clint.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("slack not working use sms")
	}
	return nil
}

// HitTheBananasError sets up payload to send error message to Hit-The-Bananas.
func HitTheBananasError(channel string, err string, logsURL string, mentions []string) error {
	errMsg := err
	ats := linkMentions(mentions)
	pld := &SendMsg{
		Text: ats + " you have a system error!",
		Attachments: []Attachments{Attachments{
			Fallback:   "error alert " + errMsg,
			Title:      "High Priority",
			AuthorName: "Hit-The_Bananas",
			AuthorIcon: "https://rpelm.com/images/banana-clipart-cartoon-1.png",
			Color:      "danger",
			Fields: []Fields{Fields{
				Title: "[Error Message]",
				Value: errMsg,
				Short: false,
			}},
			Actions: []Actions{Actions{
				Type:  "button",
				Text:  "View Logs",
				Style: "primary",
				URL:   logsURL,
			}},
		},
		},
	}
	return Send(channel, pld, nil)
}

// HitTheBananasErrorAE sets up payload to send error message to Hit-The-Bananas for app engine or any need for custom client.
func HitTheBananasErrorAE(c *http.Client, channel string, err string, logsURL string, mentions []string) error {
	errMsg := err
	ats := linkMentions(mentions)
	pld := &SendMsg{
		Text: ats + " you have a system error!",
		Attachments: []Attachments{Attachments{
			Fallback:   "error alert " + errMsg,
			Title:      "High Priority",
			AuthorName: "Hit-The_Bananas",
			AuthorIcon: "https://rpelm.com/images/banana-clipart-cartoon-1.png",
			Color:      "danger",
			Fields: []Fields{Fields{
				Title: "[Error Message]",
				Value: errMsg,
				Short: false,
			}},
			Actions: []Actions{Actions{
				Type:  "button",
				Text:  "View Logs",
				Style: "primary",
				URL:   logsURL,
			}},
		},
		},
	}
	return Send(channel, pld, c)
}

func linkMentions(ats []string) string {
	str := ""
	if len(ats) != 0 {
		for _, name := range ats {
			str += "<" + name + "> "
		}
	}
	return str
}
