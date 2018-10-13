package email

import (
	"log"
)

type nullTransport struct{}

func (n *nullTransport) SendEmail(sender, subject, body string) error {
	log.Println("email to be sent:")
	log.Println(sender, subject, body)
	return nil
}

func NewNullTransport() *nullTransport {
	return &nullTransport{}
}
