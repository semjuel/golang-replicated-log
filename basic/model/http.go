package model

type RequestMessage struct {
	Text string `json:"body"`
	W    int32  `json:"w"`
}
