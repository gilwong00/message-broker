package broker

import (
	"fmt"
	"sync"
	"time"

	"github.com/gilwong00/message-broker/pkg/message"
	"github.com/gilwong00/message-broker/pkg/subscriber"
)

// this is an in memory borker
type Broker struct {
	subscribers map[string][]*subscriber.Subscriber
	mutex       sync.Mutex
}

// NewBroker creates a new instance of the in-memory broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]*subscriber.Subscriber),
	}
}

// Subscribe allows a client to subscribe to a topic and receive messages.
func (b *Broker) Subscribe(topic string) *subscriber.Subscriber {
	// to guarantee exclusive access to the subscribers map
	// we want to lock access to each per each request. This prevents multiple
	// go routines from updating the map concurrently
	b.mutex.Lock()
	defer b.mutex.Unlock()
	subscriber := &subscriber.Subscriber{
		Channel:     make(chan interface{}, 1),
		Unsubscribe: make(chan bool),
	}
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
	return subscriber
}

// Unsubscribe removes a subscriber from a topic.
func (b *Broker) Unsubscribe(topic string, subscriber *subscriber.Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if subscribers, found := b.subscribers[topic]; found {
		for i, sub := range subscribers {
			if sub == subscriber {
				// closing channel
				// this indicates that our subscriber is no longer interested
				// in receiving messages.
				close(sub.Channel)
				// removing subscriber from slice
				b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

// Publish sends a message to all subscribers to the specifed topic.
func (b *Broker) Publish(msg *message.Message) {
	// lock the mutex to safely access the subscribers map
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if subscribers, ok := b.subscribers[msg.Topic]; ok {
		for _, sub := range subscribers {
			select {
			case sub.Channel <- msg.Payload:
				// If a subscriber is slow to receive the message (after 1 second here), consider it unresponsive and unsubscribe it from the topic.
			case <-time.After(time.Second):
				fmt.Printf("Subscriber slow. Unsubscribing from topic: %s\n", msg.Topic)
				b.Unsubscribe(msg.Topic, sub)
			}
		}
	}
}
