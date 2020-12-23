package api

import (
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"replicated_log/basic/model"
	"time"
)

const (
	SUCCESS = "success"
	FAILURE = "failure"
)

// Server represents the gRPC server.
type Server struct{}

func (s *Server) Send(ctx context.Context, in *LogMessage) (*LogMessage, error) {
	log.Printf("gRPC received message: %s \n", in.Body)

	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(20)
	log.Printf("Sleep for %d seconds \n", delay)
	time.Sleep(time.Duration(delay) * time.Second)

	model.AddMessage(model.Message{Id: in.Id, Text: in.Body})

	status := SUCCESS
	if delay%2 == 0 {
		status = FAILURE
	}
	log.Printf("Response status - `%s`", status)

	return &LogMessage{Body: status}, nil
}
