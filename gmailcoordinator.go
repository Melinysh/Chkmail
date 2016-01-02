package main

import (
	"google.golang.org/api/gmail/v1"
)

type GmailCoordinator struct {
	msgService      *gmail.UsersMessagesService
	currentMessages map[string]EmailMessage // maps message ID to email message itself
	pageToken       string
}

func NewGmailCoordinator(service *gmail.UsersMessagesService) GmailCoordinator {
	return GmailCoordinator{service, map[string]EmailMessage{}, ""}
}

func (self *GmailCoordinator) FetchMessages() []EmailMessage {
	ids := self.messagesIDs()
	return self.emailsFromIDs(ids)
}

// returns the latest list of message ids and updates self.lastPageToken as required
func (self *GmailCoordinator) messagesIDs() []string {
	listCall := self.msgService.List(*emailAddress)
	if len(self.pageToken) != 0 {
		listCall.PageToken(self.pageToken)
	}
	IDList, err := listCall.Do()
	var messageIDs []string
	if err != nil {
		debugPrint("Error fetching list of message ids", err)
		return messageIDs
	}
	self.pageToken = IDList.NextPageToken
	for _, emptyMsg := range IDList.Messages {
		messageIDs = append(messageIDs, emptyMsg.Id)
	}
	return messageIDs
}

func (self *GmailCoordinator) emailsFromIDs(ids []string) []EmailMessage {
	messages := []EmailMessage{}
	for _, id := range ids {
		rawMsg, getErr := self.msgService.Get(*emailAddress, id).Format("full").Do()
		if getErr != nil {
			debugPrint("Error fetching message with id", id, getErr)
			continue
		}
		e, parseErr := NewEmailMessage(rawMsg)
		if parseErr != nil {
			debugPrint("Error parsing email with id", id, parseErr)
			continue
		}
		messages = append(messages, e)
	}
	return messages
}
