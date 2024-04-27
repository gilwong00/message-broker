package broker

import (
	"sync"

	"github.com/gilwong00/message-broker/pkg/subscriber"
)

// this is an in memory borker
type Broker struct {
	subscribers map[string][]*subscriber.Subscriber
	mutex       sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]*subscriber.Subscriber),
	}
}
