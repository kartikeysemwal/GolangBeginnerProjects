package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/kartikeysemwal/ProtoBuff/proto"
)

func callcallSayHelloBidirectionalStreaming(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Bidirectional streaming started from server side")

	stream, err := client.SayHelloBidirectionalStreaming(context.Background())

	if err != nil {
		log.Fatalf("Not able to start stream %v", err)
	}

	first := true

	for _, name := range names.Names {
		if !first {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error in callSayHelloBidirectionalStreaming from client %v", err)
			}
			log.Printf("callSayHelloBidirectionalStreaming: Message from server %v", message.Message)
		}

		first = false

		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error in sending req from callSayHelloBidirectionalStreaming %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()

	log.Printf("Bidirectional streaming ended from server side")
}
