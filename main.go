package main

import (
	"context"
	"grpc-protobuf/pb/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func (u *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("User is created")

	return &user.CreateResponse{
		Message: "User created successfully",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal("There is error in your net listen:", err)
		return
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})

	reflection.Register(serv) // Especially used for development & testing using postman.

	if err = serv.Serve(lis); err != nil {
		log.Fatal("Error running server:", err)
	}
}

/*
	=> Reflection adalah fitur grpc yang memungkinkan kita untuk
	mengetahui layanan apa saja yang tersedia di server.
	Bisa digunakan untuk mengetahui service dan method apa saja
	yang tersedia di server (Testing via Postman dengan jenis gRPC).
	Sebaiknya digunakan di environment development, bukan production.
	Untuk production, sebaiknya menggunakan gRPC server.
*/
