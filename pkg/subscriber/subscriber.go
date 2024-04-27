package subscriber

// Subscriber belong to a channel
type Subscriber struct {
	Channel     chan interface{}
	Unsubscribe chan bool
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}
