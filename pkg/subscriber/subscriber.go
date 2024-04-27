package subscriber

import "fmt"

// Subscriber belongs to a channel
type Subscriber struct {
	// maybe give it an id
	Channel     chan interface{}
	Unsubscribe chan bool
}

func (s *Subscriber) PrintChannel() {
	for msg := range s.Channel {
		fmt.Println(">>> channel msg", msg)
	}
}

func (s *Subscriber) ListenToChannelEvents() {
	go func() {
		for {
			select {
			case msg, ok := <-s.Channel:
				if !ok {
					fmt.Println("Subscriber closed channel.")
					return
				}
				fmt.Printf("Received new msg: %v\n", msg)
			case <-s.Unsubscribe:
				fmt.Println("Unsubscribed.")
				return
			}
		}
	}()
}
