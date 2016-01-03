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

type Subscriber struct {
	emailEvents chan EmailEvent
}

func NewSubscriber() Subscriber {
	return Subscriber{make(chan EmailEvent)}
}

type Publisher struct {
	subscribers []Subscriber
}

func NewPublisher() Publisher {
	return Publisher{}
}

func NewPublisherWithSubscriber(sub Subscriber) Publisher {
	pub := NewPublisher()
	pub.AddSubscriber(sub)
	return pub
}

func (self *Publisher) AddSubscriber(sub Subscriber) {
	self.subscribers = append(self.subscribers, sub)
}

func (self *Publisher) RemoveSubscriber(sub Subscriber) {
	var newSubs []Subscriber
	// filter out all other subs into new array
	for _, s := range self.subscribers {
		if s != sub {
			newSubs = append(newSubs, s)
		}
	}
	self.subscribers = newSubs
}

func (self *Publisher) Publish(event EmailEvent) {
	for _, sub := range self.subscribers {
		sub.emailEvents <- event
	}
}
