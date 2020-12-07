package api

import (
	"golang.org/x/net/context"
	"log"
	"replicated_log/basic/model"
)

// Server represents the gRPC server.
type Server struct{}

func (s *Server) Send(ctx context.Context, in *LogMessage) (*LogMessage, error) {
	log.Printf("gRPC received message: %s", in.Body)

	model.AddMessage(model.Message{Text: in.Body})

	return &LogMessage{Body: "ok"}, nil
}
