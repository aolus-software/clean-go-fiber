package mailer

import (
	"strconv"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg *config.MailConfig
}

func New(cfg *config.MailConfig) *Mailer {
	return &Mailer{cfg: cfg}
}

type Message struct {
	To      string
	Subject string
	Body    string
}

func (m *Mailer) Send(msg Message) error {
	port, err := strconv.Atoi(m.cfg.Port)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", m.cfg.MailFrom)
	mail.SetHeader("To", msg.To)
	mail.SetHeader("Subject", msg.Subject)
	mail.SetBody("text/html", msg.Body)

	dialer := gomail.NewDialer(m.cfg.Host, port, m.cfg.User, m.cfg.Password)
	return dialer.DialAndSend(mail)
}
