package main

import (
	"context"
	"grpc-protobuf/pb/chat"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}

	chatClient := chat.NewChatServiceClient(clientConn)

	stream, err := chatClient.SendMessage(context.Background())

	if err != nil {
		log.Fatal("Failed to send message: ", err)
	}

	// Looping sending messages (Client Streaming) 10x (Ex: user send message)
	for i := 1; i <= 10; i++ {
		err = stream.Send(&chat.ChatMessage{
			UserId:  strconv.Itoa(i),
			Content: "Hello from client",
		})

		if err != nil {
			log.Fatal("Failed to send message: ", err)
		}

	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal("Failed to close and receive message: ", err)
	}

	log.Println("Connection is closed. Message: ", res.Message)
}
