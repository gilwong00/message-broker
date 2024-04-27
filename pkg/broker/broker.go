package broker

import (
	"fmt"
	"sync"
	"time"

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

func (b *Broker) Unsubscribe(topic string, subscriber *subscriber.Subscriber) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if subscribers, ok := b.subscribers[topic]; ok {
		for i, sub := range subscribers {
			if sub == subscriber {
				// closing channel
				// this indicates that our subscriber is no longer interested
				// in receiving messages.
				close(sub.Channel)
				// removing subscriber from slice
				b.subscribers[topic] = append(subscribers[:i], subscribers[:i+1]...)
				return nil
			}
		}
		return nil
	}
	return fmt.Errorf("topic does not exist, %s", topic)
}

func (b *Broker) Publish(topic string, message interface{}) {
	// lock the mutex to safely access the subscribers map
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if subscribers, ok := b.subscribers[topic]; ok {
		for _, sub := range subscribers {
			select {
			case sub.Channel <- message:
			// If a subscriber is slow to receive the message (after 1 second here), consider it unresponsive and unsubscribe it from the topic.
			case <-time.After(time.Second):
				fmt.Printf("Slow subscription. Unsubscribing from topic: %s\n", topic)
				b.Unsubscribe(topic, sub)
			}
		}
	}
}
