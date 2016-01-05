package main

import (
	"google.golang.org/api/gmail/v1"
)

type GmailCoordinator struct {
	EmailPublisher
	sub             UISubscriber
	service         *gmail.Service
	msgService      *gmail.UsersMessagesService
	draftService    *gmail.UsersDraftsService
	labelService    *gmail.UsersLabelsService
	attachService   *gmail.UsersMessagesAttachmentsService
	threadService   *gmail.UsersThreadsService
	currentMessages map[string]EmailMessage // maps message ID to email message itself
	pageToken       string
}

func NewGmailCoordinator(service *gmail.Service) GmailCoordinator {
	return GmailCoordinator{
		NewEmailPublisher(),
		NewUISubscriber(),
		service,
		gmail.NewUsersMessagesService(service),
		gmail.NewUsersDraftsService(service),
		gmail.NewUsersLabelsService(service),
		gmail.NewUsersMessagesAttachmentsService(service),
		gmail.NewUsersThreadsService(service),
		map[string]EmailMessage{},
		""}

}

func NewGmailCoordinatorWithSubscriber(service *gmail.Service, sub EmailSubscriber) GmailCoordinator {
	gc := NewGmailCoordinator(service)
	gc.EmailPublisher.AddSubscriber(sub)
	return gc
}

func (self *GmailCoordinator) ListenForUIChanges() {
	go func() {
		for {
			event := <-self.sub.uiEvents
			switch event.Action {
			case Trash:
				self.trashMessage(event.Email.Id)
			default:
				debugPrint("Unimplemented for action", event.Action)
			}
		}

	}()
}

func (self *GmailCoordinator) FetchMessages() {
	ids := self.messagesIDs()
	self.emailsFromIDs(ids)
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

func (self *GmailCoordinator) emailsFromIDs(ids []string) {
	for _, id := range ids {
		rawMsg, getErr := self.msgService.Get(*emailAddress, id).Do()
		if getErr != nil {
			debugPrint("Error fetching message with id", id, getErr)
			continue
		}
		e, parseErr := NewEmailMessage(rawMsg)
		if parseErr != nil {
			debugPrint("Error parsing email with id", id, parseErr)
			continue
		}
		self.Publish(EmailEvent{Recieved, e})
	}
}

func (self *GmailCoordinator) trashMessage(id string) {
	msg, err := self.msgService.Trash(*emailAddress, id).Do()
	if err != nil {
		debugPrint("Unable to trash email with id", id, err)
		return
	}
	email, parseErr := NewEmailMessage(msg)
	if parseErr != nil {
		debugPrint("Unable to parse newly trashed email", parseErr)
		return
	}
	self.Publish(EmailEvent{Trashed, email})
}
