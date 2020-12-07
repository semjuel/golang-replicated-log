package model

var messages []Message

type Message struct {
	Text string `json:"body"`
}

func AddMessage(message Message) {
	messages = append(messages, message)
}

func GetMessages() []Message {
	if len(messages) == 0 {
		messages = make([]Message, 0)
	}
	return messages
}
