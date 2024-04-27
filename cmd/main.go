package main

import (
	"time"

	"github.com/gilwong00/message-broker/pkg/broker"
	"github.com/gilwong00/message-broker/pkg/message"
)

const testTopic = "test_topic"

func main() {
	broker := broker.NewBroker()
	subscriber := broker.Subscribe(testTopic)
	subscriber.ListenToChannelEvents()

	broker.Publish(&message.Message{
		Topic:   testTopic,
		Payload: "Hello msg",
	})
	broker.Publish(&message.Message{
		Topic:   testTopic,
		Payload: "This is another test message",
	})
	time.Sleep(2 * time.Second)
	broker.Unsubscribe(testTopic, subscriber)
	broker.Publish(&message.Message{
		Topic:   "Unsuccessful topic",
		Payload: "This is wont print",
	})
	time.Sleep(time.Second)
}
