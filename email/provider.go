package email

type Provider interface {
	SendEmail(sender, subject, body string) error
}
