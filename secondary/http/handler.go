package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"replicated_log/basic/model"
)

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(w, r)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(model.GetMessages())
	if err != nil {
		log.Printf("Encoder error %v \n", err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(w, r)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validate(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	if r.Method == "OPTIONS" {
		return errors.New("skip options")
	}
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		log.Println(fmt.Sprintf("Request method %s not allowed", r.Method))
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		return errors.New("method not allowed")
	}

	return nil
}
