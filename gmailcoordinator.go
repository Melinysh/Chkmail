package main

import (
	"google.golang.org/api/gmail/v1"
)

type GmailCoordinator struct {
	msgService *gmail.UsersMessagesService
}

func NewGmailCoordinator(service *gmail.UsersMessagesService) GmailCoordinator {
	return GmailCoordinator{service}
}
