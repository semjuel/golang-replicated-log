package model

import (
	"sort"
	"sync/atomic"
)

type Messages []Message

type Message struct {
	Id   int32  `json:"id"`
	Text string `json:"body"`
}

func (a Messages) Len() int           { return len(a) }
func (a Messages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Messages) Less(i, j int) bool { return a[i].Id < a[j].Id }

func (m Message) intId() int { return int(m.Id) }

var id int32 = 0
var messages Messages

func InitMessage(text string) Message {
	message := Message{
		Id:   atomic.LoadInt32(&id),
		Text: text,
	}
	atomic.AddInt32(&id, 1)

	return message
}

func AddMessage(message Message) {
	if len(messages) == 0 {
		messages = make([]Message, 0)
	}

	exists := false
	for _, v := range messages {
		if v.Id == message.Id {
			exists = true
			break
		}
	}

	if exists == false {
		messages = append(messages, message)
	}
}

func GetMessages() Messages {
	var response = make(Messages, 0)

	if len(messages) == 0 {
		return response
	}

	sort.Sort(messages)
	first := messages[0]

	if first.intId() != 0 {
		return response
	}

	response = append(response, first)
	for i := 1; i < len(messages); i++ {
		message := messages[i]
		if message.intId() != (response[i-1].intId() + 1) {
			return response
		}

		response = append(response, message)
	}

	return response
}
