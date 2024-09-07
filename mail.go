package main

import (
	"github.com/wneessen/go-mail"
	"google-contacts-birthday-notification/config"
)

type MailClientInterface interface {
	DialAndSend(messages ...*mail.Msg) error
}

type MailService struct {
	client MailClientInterface
	config *config.Config
}

func NewMailService(config *config.Config) *MailService {
	var tlsPolicy mail.TLSPolicy
	if config.Mail.Tls && config.Mail.Secure {
		tlsPolicy = mail.TLSMandatory
	} else if config.Mail.Tls && !config.Mail.Secure {
		tlsPolicy = mail.TLSOpportunistic
	} else {
		tlsPolicy = mail.NoTLS
	}
	client, _ := mail.NewClient(
		config.Mail.Host,
		mail.WithPort(config.Mail.Port),
		mail.WithUsername(config.Mail.User),
		mail.WithPassword(config.Mail.Password),
		mail.WithSSL(),
		mail.WithTLSPolicy(tlsPolicy),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	return &MailService{config: config, client: client}
}

func (c *MailService) SendMail(message string) error {
	m := mail.NewMsg()
	if err := m.From(c.config.Mail.Sender); err != nil {
		return err
	}
	if err := m.To(c.config.Mail.Receiver); err != nil {
		return err
	}
	m.Subject("Birthday Notification")
	m.SetBodyString("text/plain", message)

	return c.client.DialAndSend(m)
}
