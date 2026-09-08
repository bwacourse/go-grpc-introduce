package main

import (
	"context"
	"grpc-protobuf/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to create client: ", err)
	}

	userClient := user.NewUserServiceClient(clientConn)

	response, err := userClient.CreateUser(context.Background(), &user.User{
		Id:      1,
		Age:     23,
		Balance: 130000,
		Address: &user.Address{
			Id:          123,
			FullAddress: "Jl. Coding No. 123",
			Province:    "Jawa Barat",
			City:        "Jakarta",
		},
	})

	if err != nil {
		log.Fatal("Error calling user client: ", err)
	}

	log.Println("Got message from server: ", response.Message)
}
