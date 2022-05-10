package stream

import (
	"log"
	"sync"
	"time"
)

type subscriber struct {
	wg      *sync.WaitGroup
	channel Delegate
}

func (s *subscriber) shuffle(item Item) {
	s.wg.Wait()
	log.Printf("Attempt to shuffle data : %d to %s ", item.Id, s.channel.name)
	select {
	case s.channel.channel <- item:
		log.Printf("Successfully shuffled data : %d to %s ", item.Id, s.channel.name)
	case <-time.After(Delay):
		log.Printf("Failed to shuffle data : %d to %s ", item.Id, s.channel.name)
	}

}
