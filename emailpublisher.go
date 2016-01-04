package main

type EmailAction int

const (
	Sent     EmailAction = iota
	Recieved             = iota
	Trashed              = iota
	Drafted              = iota
	Modified             = iota
)

type EmailEvent struct {
	Action EmailAction
	Email  EmailMessage
}

type EmailSubscriber struct {
	emailEvents chan EmailEvent
}

func NewEmailSubscriber() EmailSubscriber {
	return EmailSubscriber{make(chan EmailEvent)}
}

type EmailPublisher struct {
	subscribers []EmailSubscriber
}

func NewEmailPublisher() EmailPublisher {
	return EmailPublisher{}
}

func NewEmailPublisherWithSubscriber(sub EmailSubscriber) EmailPublisher {
	pub := NewEmailPublisher()
	pub.AddSubscriber(sub)
	return pub
}

func (self *EmailPublisher) AddSubscriber(sub EmailSubscriber) {
	self.subscribers = append(self.subscribers, sub)
}

func (self *EmailPublisher) RemoveSubscriber(sub EmailSubscriber) {
	var newSubs []EmailSubscriber
	// filter out all other subs into new array
	for _, s := range self.subscribers {
		if s != sub {
			newSubs = append(newSubs, s)
		}
	}
	self.subscribers = newSubs
}

func (self *EmailPublisher) Publish(event EmailEvent) {
	for _, sub := range self.subscribers {
		sub.emailEvents <- event
	}
}
