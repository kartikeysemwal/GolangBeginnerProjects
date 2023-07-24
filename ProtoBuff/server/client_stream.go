package main

import (
	"io"
	"log"

	pb "github.com/kartikeysemwal/ProtoBuff/proto"
)

func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	log.Printf("Streaming started for callSayHelloClientStreaming")
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming from server side %v", err)
		}
		log.Println(message)
	}

	log.Printf("Streaming ended for callSayHelloClientStreaming")

	return nil
}
