package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"time"
)

type EmailMessage struct {
	Id      string
	Body    string
	To      string
	From    string
	Subject string
	Date    time.Time
}

func NewEmailMessage(m *gmail.Message) (EmailMessage, error) {
	id := m.Id
	if m.Payload == nil {
		return EmailMessage{}, errors.New("Email does not have a payload")
	}
	body := decode(bodyTextForGmailMessage(m.Payload))
	to := ""
	from := ""
	subj := ""
	for _, h := range m.Payload.Headers {
		switch h.Name {
		case "To":
			to = h.Value
			continue
		case "From":
			from = h.Value
			continue
		case "Subject":
			subj = h.Value
			continue
		}
	}
	date := time.Unix(m.InternalDate/1000, 0)
	return EmailMessage{id, body, to, from, subj, date}, nil
}

func (self EmailMessage) Print() {
	fmt.Println("To:", self.To)
	fmt.Println("From:", self.From)
	fmt.Println(self.Date.Format("Mon Jan _2 15:04:05 2006"))
	fmt.Println("Subject:", self.Subject)
	fmt.Println("\n", self.Body)
}

func bodyTextForGmailMessage(m *gmail.MessagePart) string {
	body := ""
	if m.MimeType == "text/plain" {
		body = m.Body.Data
	}
	for _, p := range m.Parts {
		body += bodyTextForGmailMessage(p)
	}
	return body
}

func decode(str string) string {
	data, _ := base64.StdEncoding.DecodeString(str)
	/*if err != nil {
		panic(err)
	}*/
	return string(data)
}
