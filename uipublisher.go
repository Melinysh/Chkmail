package main

type UIAction int

const (
	Trash    UIAction = iota
	Reply             = iota
	ReplyAll          = iota
	Forward           = iota
	Draft             = iota
)

type UIEvent struct {
	Action UIAction
	Email  EmailMessage
}

type UISubscriber struct {
	uiEvents chan UIEvent
}

func NewUISubscriber() UISubscriber {
	return UISubscriber{make(chan UIEvent)}
}

type UIPublisher struct {
	subscribers []UISubscriber
}

func NewUIPublisher() UIPublisher {
	return UIPublisher{}
}

func NewUIPublisherWithSubscriber(sub UISubscriber) UIPublisher {
	pub := NewUIPublisher()
	pub.AddSubscriber(sub)
	return pub
}

func (self *UIPublisher) AddSubscriber(sub UISubscriber) {
	self.subscribers = append(self.subscribers, sub)
}

func (self *UIPublisher) RemoveSubscriber(sub UISubscriber) {
	var newSubs []UISubscriber
	// filter out all other subs into new array
	for _, s := range self.subscribers {
		if s != sub {
			newSubs = append(newSubs, s)
		}
	}
	self.subscribers = newSubs
}

func (self *UIPublisher) Publish(event UIEvent) {
	for _, sub := range self.subscribers {
		sub.uiEvents <- event
	}
}
