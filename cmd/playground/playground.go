package main

import (
	"log"
	"time"
)

func main() {

	c1 := make(chan string)

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				log.Printf("Boom Boom Boom")
			}
		}
	}()

	go func() {
		select {
		case c1 <- "try":
			log.Printf("Got one from channel Ok")
		case <-time.After(2 * time.Second):
			log.Printf("Timeout")
		}
	}()

	time.Sleep(4 * time.Second)

	<-c1

	time.Sleep(10 * time.Second)

}
