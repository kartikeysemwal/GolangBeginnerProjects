package main

import (
	"context"
	"log"
	"time"

	pb "github.com/kartikeysemwal/ProtoBuff/proto"
)

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Streaming started from client for SayHelloClientStreaming")

	stream, err := client.SayHelloClientStreaming(context.Background())

	if err != nil {
		log.Fatalf("Not able to start stream %v", err)
	}

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error in sending req from callSayHelloClientStreaming %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()

	log.Printf("Streaming ended")
}
