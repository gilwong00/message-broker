package subscriber

// Subscriber belong to a channel
type Subscriber struct {
	// maybe give it an id
	Channel     chan interface{}
	Unsubscribe chan bool
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}
