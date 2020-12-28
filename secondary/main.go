package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"replicated_log/basic/api"
	"replicated_log/basic/model"
	"replicated_log/basic/server"
	"replicated_log/secondary/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	name := os.Getenv("REPLICATED_LOG_NODE_NAME")
	if name == "" {
		name = "Secondary"
	}
	log.Println(fmt.Sprintf("%s started...", name))

	routes := model.Routes{
		{Pattern: "/messages", Handler: http.MessagesHandler},
		{Pattern: "/health", Handler: http.HealthHandler},
	}
	go server.Run(name, routes)

	gRPCPort := os.Getenv("REPLICATED_LOG_GRPC_PORT")
	log.Printf("Start gRPC server on port %s", gRPCPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := api.Server{}
	grpcServer := grpc.NewServer()
	api.RegisterLogMessageServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
