package main

import (
	"io"
	"log"

	pb "github.com/kartikeysemwal/ProtoBuff/proto"
)

func (s *helloServer) SayHelloBidirectionalStreaming(stream pb.GreetService_SayHelloBidirectionalStreamingServer) error {
	log.Printf("Bidirectional streaming started from server side")

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error in SayHelloBidirectionalStreaming from server %v", err)
		}
		log.Printf("SayHelloBidirectionalStreaming: Message from client %v", message.Name)
		res := &pb.HelloResponse{
			Message: "Hello " + message.Name + ", from server side",
		}
		if err := stream.Send(res); err != nil {
			log.Fatalf("Error while sending response from server for SayHelloBidirectionalStreaming %v", err)
		}
	}

	log.Printf("Bidirectional streaming ended from server side")
	return nil
}
