package main

import (
	"context"
	"errors"
	"grpc-protobuf/pb/chat"
	"grpc-protobuf/pb/user"
	"io"
	"log"
	"net"

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

		log.Printf("Receive Message %s, to %s", req.Content, req.UserId)
	}

	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Thanks for the message!",
	})
}

// func (c chatService) ReceiveMessage(*chat.ReceiveMessageRequest, grpc.ServerStreamingServer[chat.ChatResponse]) error {
// 	return status.Error(codes.Unimplemented, "method ReceiveMessage not implemented")
// }
// func (c chatService) Chat(grpc.BidiStreamingServer[ChatMessage, ChatMessage]) error {
// 	return status.Error(codes.Unimplemented, "method Chat not implemented")
// }

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
