package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"replicated_log/basic/api"
	"replicated_log/basic/model"
	"replicated_log/basic/server"
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

	routes := model.Routes{{Pattern: "/messages", Handler: handler}}
	go server.Run(name, routes)

	log.Println("Start gRPC server")
	gRPCPort := os.Getenv("REPLICATED_LOG_GRPC_PORT")
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

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		log.Println(fmt.Sprintf("Request method %s not allowed", r.Method))
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(model.GetMessages())
	if err != nil {
		log.Printf("Encoder error %v \n", err)
	}
}
