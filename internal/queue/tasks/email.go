package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailSend     = "email:send"
	EmailSendMaxRetry = 3
)

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
	return asynq.NewTask(TypeEmailSend, data, asynq.MaxRetry(EmailSendMaxRetry)), nil
}

func EnqueueEmail(ctx context.Context, client *asynq.Client, payload EmailSendPayload) error {
	task, err := NewEmailSendTask(payload)
	if err != nil {
		return err
	}
	_, err = client.EnqueueContext(ctx, task)
	return err
}
