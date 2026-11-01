package mailer

import (
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
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.cfg.From)
	mail.SetHeader("To", msg.To)
	mail.SetHeader("Subject", msg.Subject)
	mail.SetBody("text/html", msg.Body)

	dialer := gomail.NewDialer(m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Pass)
	dialer.SSL = m.cfg.Secure
	return dialer.DialAndSend(mail)
}
