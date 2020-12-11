package api

import (
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"replicated_log/basic/model"
	"time"
)

// Server represents the gRPC server.
type Server struct{}

func (s *Server) Send(ctx context.Context, in *LogMessage) (*LogMessage, error) {
	log.Printf("gRPC received message: %s \n", in.Body)

	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(20)
	log.Printf("Sleep for %d seconds \n", delay)
	time.Sleep(time.Duration(delay) * time.Second)

	model.AddMessage(model.Message{Text: in.Body})

	return &LogMessage{Body: "ok"}, nil
}
