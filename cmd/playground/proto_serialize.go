package main

import (
	"google.golang.org/protobuf/proto"
	"log"
	"pilsner/internal/communication"
	"pilsner/proto/pb"
	"pilsner/translator"
)

func main() {

	item := communication.Item{
		Id:      1543535,
		Content: "dummy",
		Source:  "playground",
	}

	_, itemPb := translator.Translate[pb.Item](&item)

	data, _ := proto.Marshal(&itemPb)

	log.Printf("marshaled %d data\n", len(data))
}
