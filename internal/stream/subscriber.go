package stream

import (
	"log"
	"sync"
)

type subscriber struct {
	wg      *sync.WaitGroup
	channel Delegate
}

func (s *subscriber) shuffle(item Item) {
	s.wg.Wait()
	log.Printf("Attempt to shuffle data : %d to %s ", item.Id, s.channel.name)
	s.channel.channel <- item
	log.Printf("Successfully shuffled data : %d to %s ", item.Id, s.channel.name)
}
