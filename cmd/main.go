package main

import (
	"github.com/gilwong00/message-broker/pkg/broker"
)

func main() {
	broker := broker.NewBroker()

	subscriber := broker.Subscribe("test_topic")
	broker.Unsubscribe("test_topic", subscriber)
}
