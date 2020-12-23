package http

import (
	"encoding/json"
	"log"
	"net/http"
	"replicated_log/basic/model"
	"replicated_log/master/grpc"
	"replicated_log/master/service"
	"replicated_log/master/utils"
	"sync/atomic"
)

var concern int32

func Handler(w http.ResponseWriter, r *http.Request) {
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
		var m model.RequestMessage
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

		message := model.InitMessage(m.Text)
		model.AddMessage(message)

		var iterations int32 = 1
		for _, v := range service.GetSecondaries() {
			go grpc.Replicate(message, v.GRPS, &iterations)
		}

		concern = 1
		if m.W > 1 {
			concern = m.W
		} else if concern > 3 {
			concern = 3
		}

		for {
			state := atomic.LoadInt32(&iterations)
			if state >= concern {
				break
			}
		}

		w.WriteHeader(http.StatusCreated)
		break

	default:
		log.Printf("Request method %s not allowed \n", r.Method)
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		break
	}
}
