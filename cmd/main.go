package main

import (
	"fmt"
	"log"
	"my_module/api"
	"my_module/api/handler"
	"my_module/generated/auth_service"
	"my_module/logs"
	"my_module/service"
	"my_module/storage/postgres"
	"net"

	"google.golang.org/grpc"
)

func main() {
	logs.InitLogger()
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("error while listening: %v", err)
		return
	}
	defer listener.Close()

	userService := service.NewUserService(db)
	server := grpc.NewServer()
	auth_service.RegisterAuthServiceServer(server, userService)

	log.Printf("server listening at %v", listener.Addr())

	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Fatalf("error while serving: %v", err)
		}
	}()

	router := api.Router(handler.NewHandler(postgres.NewUserRepo(db), logs.Logger))
	log.Println("server is running")
	log.Fatal(router.Run(":8085"))
}
