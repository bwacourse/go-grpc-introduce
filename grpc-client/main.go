package main

import (
	"context"
	"errors"
	"grpc-protobuf/pb/chat"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* [ Client Streaming ] */
// func main() {
// 	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	if err != nil {
// 		log.Fatal("Failed to create client: ", err)
// 	}

// 	chatClient := chat.NewChatServiceClient(clientConn)

// 	stream, err := chatClient.SendMessage(context.Background())

// 	if err != nil {
// 		log.Fatal("Failed to send message: ", err)
// 	}

// 	// Looping sending messages (Client Streaming) 10x (Ex: user send message)
// 	for i := 1; i <= 10; i++ {
// 		err = stream.Send(&chat.ChatMessage{
// 			UserId:  int64(i),
// 			Content: "Hello from client",
// 		})

// 		if err != nil {
// 			log.Fatal("Failed to send message: ", err)
// 		}

// 	}

// 	res, err := stream.CloseAndRecv()

// 	if err != nil {
// 		log.Fatal("Failed to close and receive message: ", err)
// 	}

// 	log.Println("Connection is closed. Message: ", res.Message)
// }

/* [ Server Streaming ] */
func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}

	chatClient := chat.NewChatServiceClient(clientConn)

	stream, err := chatClient.ReceiveMessage(context.Background(), &chat.ReceiveMessageRequest{
		UserId: 123,
	})

	// Looping menerima message (Server Streaming) tanpa batas (Infinity Loop)
	for {
		msg, err := stream.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal("Failed to send message: ", err)
		}

		log.Printf("Got message to %d, content %s", msg.UserId, msg.Content)
	}
}
