package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"replicated_log/basic/model"
	"replicated_log/basic/server"
	"replicated_log/master/http"
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
	routes := model.Routes{{Pattern: "/messages", Handler: http.Handler}}
	server.Run(name, routes)
}
