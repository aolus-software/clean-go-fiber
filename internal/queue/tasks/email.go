package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeEmailSend = "email:send"

type EmailSendPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEmailSendTask(payload EmailSendPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailSend, data), nil
}
