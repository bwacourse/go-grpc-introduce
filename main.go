package main

import (
	"context"
	"errors"
	"grpc-protobuf/pb/chat"
	"grpc-protobuf/pb/user"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

type chatService struct {
	chat.UnimplementedChatServiceServer
}

func (u *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("User is created")

	return &user.CreateResponse{
		Message: "User created successfully",
	}, nil
}

func (c chatService) SendMessage(stream grpc.ClientStreamingServer[chat.ChatMessage, chat.ChatResponse]) error {
	// Looping untuk menerima message (Client Streaming) tanpa batas (Infinity Loop)
	// Selama belum ada error atau EOF (End Of File), maka akan terus berjalan.
	for {
		req, err := stream.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "Error receiving message: %v", err)
		}

		log.Printf("Receive Message %s, to %d", req.Content, req.UserId)
	}

	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Thanks for the message!",
	})
}

func (c chatService) ReceiveMessage(request *chat.ReceiveMessageRequest, stream grpc.ServerStreamingServer[chat.ChatMessage]) error {
	log.Printf("Got connection request from %d\n", request.UserId)

	// Looping untuk mengirim message (Server Streaming) 10x
	for i := 1; i <= 10; i++ {
		err := stream.Send(&chat.ChatMessage{
			UserId:  123,
			Content: "Hi",
		})

		if err != nil {
			return status.Errorf(codes.Unknown, "Error sending message to client: %v", err)
		}
	}

	return nil
}

// Bidirectional Streaming RPC
func (c chatService) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	for {
		msg, err := stream.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "error receiving message")
		}

		log.Printf("Got message form %d, content: %s", msg.UserId, msg.Content)

		time.Sleep(2 * time.Second) // Waiting for 2 seconds

		err = stream.Send(&chat.ChatMessage{
			UserId:  50,
			Content: "Reply from server",
		})

		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message to client: %v", err)
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal("There is error in your net listen:", err)
		return
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv) // Especially used for development & testing using postman.

	if err = serv.Serve(lis); err != nil {
		log.Fatal("Error running server:", err)
	}
}
