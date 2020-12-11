package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
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
		name = "Master"
	}
	log.Println(fmt.Sprintf("%s started...", name))

	// Run HTTP server.
	routes := model.Routes{{Pattern: "/messages", Handler: handler}}
	server.Run(name, routes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(model.GetMessages())
		if err != nil {
			log.Printf("Encoder error %v \n", err)
		}
		break

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var m model.Message
		err := decoder.Decode(&m)
		if err != nil {
			log.Printf("Decode error %s \n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if m.Text == "" {
			log.Println("Empty message body")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		log.Printf("Message recieved: %s", m.Text)

		model.AddMessage(m)

		for i := 1; i <= 2; i++ {
			target := fmt.Sprintf("replicated-log-secondary-%d:800%d", i, i)
			sendMessageToSecondary(m.Text, target)
		}

		w.WriteHeader(http.StatusCreated)
		break

	default:
		log.Printf("Request method %s not allowed \n", r.Method)
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		break
	}
}

func sendMessageToSecondary(body string, target string) {
	// Run gRPC client.
	log.Printf("Start gRPC client - %s \n", target)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Printf("Did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewLogMessageServiceClient(conn)

	response, err := c.Send(context.Background(), &api.LogMessage{Body: body})
	if err != nil {
		log.Printf("Server: %s. Error when calling Send: %s", target, err)
	}

	log.Printf("Response from server %s: %s", target, response.Body)
}
