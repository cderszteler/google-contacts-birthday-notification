package main

import (
	"github.com/wneessen/go-mail"
)

var client MailClientInterface

type MailClientInterface interface {
	DialAndSend(messages ...*mail.Msg) error
}

func createMailClient() {
	var tlsPolicy mail.TLSPolicy
	if config.Mail.Tls && config.Mail.Secure {
		tlsPolicy = mail.TLSMandatory
	} else if config.Mail.Tls && !config.Mail.Secure {
		tlsPolicy = mail.TLSOpportunistic
	} else {
		tlsPolicy = mail.NoTLS
	}
	client, _ = mail.NewClient(
		config.Mail.Host,
		mail.WithPort(config.Mail.Port),
		mail.WithUsername(config.Mail.User),
		mail.WithPassword(config.Mail.Password),
		mail.WithSSL(),
		mail.WithTLSPolicy(tlsPolicy),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
}

func SendMail(message string) error {
	if client == nil {
		createMailClient()
	}

	m := mail.NewMsg()
	if err := m.From(config.Mail.Sender); err != nil {
		return err
	}
	if err := m.To(config.Mail.Receiver); err != nil {
		return err
	}
	m.Subject("Birthday Notification")
	m.SetBodyString("text/plain", message)

	return client.DialAndSend(m)
}
