package email

import (
	"github.com/mailgun/mailgun-go"
)

type mailgunTransport struct {
	client      mailgun.Mailgun
	targetEmail string
	domain      string
	secretKey   string
}

func (m *mailgunTransport) SendEmail(sender, subject, body string) error {
	message := m.client.NewMessage(sender, subject, body, m.targetEmail)
	if _, _, err := m.client.Send(message); err != nil {
		return err
	}
	return nil
}

func NewMailGunTransport(targetEmail string, client mailgun.Mailgun) *mailgunTransport {
	return &mailgunTransport{
		client:      client,
		targetEmail: targetEmail,
	}
}
