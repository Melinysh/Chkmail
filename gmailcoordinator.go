package main

import (
	"google.golang.org/api/gmail/v1"
)

type GmailCoordinator struct {
	msgService      *gmail.UsersMessagesService
	currentMessages map[string]EmailMessage // maps message ID to email message itself
}

func NewGmailCoordinator(service *gmail.UsersMessagesService) GmailCoordinator {
	return GmailCoordinator{service, map[string]EmailMessage{}}
}

func (self GmailCoordinator) FetchLatestMessages() []EmailMessage {
	IDList, err := self.msgService.List(*emailAddress).Do()
	var messages []EmailMessage
	if err != nil {
		debugPrint("Error fetching list of message ids", err)
		return messages
	}
	for _, emptyMsg := range IDList.Messages {
		rawMsg, getErr := self.msgService.Get(*emailAddress, emptyMsg.Id).Format("full").Do()
		if getErr != nil {
			debugPrint("Error fetching message with id", emptyMsg.Id, getErr)
			continue
		}
		e, parseErr := NewEmailMessage(rawMsg)
		if parseErr != nil {
			debugPrint("Error parsing email with id", emptyMsg.Id, parseErr)
			continue
		}
		messages = append(messages, e)
	}
	return messages
}
