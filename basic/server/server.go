package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"replicated_log/basic/model"
)

const (
	defaultPort = "7085"
)

func Run(name string, routes model.Routes) {
	port := os.Getenv("REPLICATED_LOG_HTTP_PORT")
	if port == "" {
		port = defaultPort
	}
	port = fmt.Sprintf(":%s", port)

	for _, r := range routes {
		log.Println(fmt.Sprintf("Bind route - %s", r.Pattern))
		http.HandleFunc(r.Pattern, r.Handler)
	}

	log.Println(fmt.Sprintf("Start %s http server at %s port", name, port))
	log.Fatal(http.ListenAndServe(port, nil))
}
