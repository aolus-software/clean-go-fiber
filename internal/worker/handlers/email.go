package handlers

import (
	"context"
	"encoding/json"

	"github.com/aolus-software/clean-go-fiber/internal/libs/mailer"
	"github.com/aolus-software/clean-go-fiber/internal/queue/tasks"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type EmailHandler struct {
	mailer *mailer.Mailer
	log    *zap.Logger
}

func NewEmailHandler(m *mailer.Mailer, log *zap.Logger) *EmailHandler {
	return &EmailHandler{mailer: m, log: log}
}

func (h *EmailHandler) HandleEmailSend(ctx context.Context, t *asynq.Task) error {
	var payload tasks.EmailSendPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	if err := h.mailer.Send(mailer.Message{
		To:      payload.To,
		Subject: payload.Subject,
		Body:    payload.Body,
	}); err != nil {
		h.log.Error("failed to send email", zap.String("to", payload.To), zap.Error(err))
		return err
	}

	h.log.Info("email sent", zap.String("to", payload.To))
	return nil
}
