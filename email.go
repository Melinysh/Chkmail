package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jaytaylor/html2text"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"mime/quotedprintable"
	"strings"
	"time"
)

type EmailMessage struct {
	Id      string
	Body    string
	To      string
	From    string
	Subject string
	Date    time.Time
	Labels  []string
}

func NewEmailMessage(m *gmail.Message) (EmailMessage, error) {
	id := m.Id
	if m.Payload == nil {
		return EmailMessage{}, errors.New("Email does not have a payload")
	}
	body := bodyTextForGmailMessage(m.Payload)
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
	labels := m.LabelIds
	return EmailMessage{id, body, to, from, subj, date, labels}, nil
}

func (self EmailMessage) Print() {
	fmt.Println("To:", self.To)
	fmt.Println("From:", self.From)
	fmt.Println(self.Date.Format("Mon Jan _2 15:04:05 2006"))
	fmt.Println("Subject:", self.Subject)
	fmt.Print("Labels: ")
	sep := ""
	for _, l := range self.Labels {
		fmt.Print(sep, l)
		sep = ", "
	}
	fmt.Println()
	fmt.Println(self.Body)
}

func (self EmailMessage) ToString() string {
	msg := "To: " + self.To + "\n"
	msg += "From: " + self.From + "\n"
	msg += self.Date.Format("Mon Jan _2 15:04:05 2006") + "\n"
	msg += "Subject: " + self.Subject + "\n"
	if *debug {
		msg += "ID: " + self.Id + "\n"
	}
	msg += "Labels: "
	sep := ""
	for _, l := range self.Labels {
		msg += sep + l
		sep = ", "
	}
	msg += "\n" + self.Body

	return msg
}

func bodyTextForGmailMessage(m *gmail.MessagePart) string {
	body := ""
	if m.MimeType == "text/plain" || m.MimeType == "text/html" {
		body = decodeBase64(m.Body.Data)
		if m.MimeType == "text/html" {
			body, _ = html2text.FromString(body)
		}
	}
	for _, p := range m.Parts {
		body += bodyTextForGmailMessage(p)
	}
	//	body = strings.Replace(body, "\n", "", -1) // fixes weird spacing in main window
	return body
}

func decodeQuotePrintable(str string) string {
	reader := strings.NewReader(str)
	qR := quotedprintable.NewReader(reader)
	decoded, err := ioutil.ReadAll(qR)
	if err != nil {
		debugPrint("Error decoding quote printable string", err)
		return ""
	}
	debugPrint(string(decoded))
	return string(decoded[:])
}

func decodeBase64(str string) string {
	fmtStr := strings.Replace(strings.Replace(str, "-", "+", -1), "_", "/", -1)
	data, err := base64.StdEncoding.DecodeString(fmtStr)
	if err != nil {
		debugPrint("Error decoding base64", err)
		return str
	}
	return string(data)
}
