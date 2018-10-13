package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mailgun/mailgun-go"

	"github.com/cjsaylor/contact-api/config"
	"github.com/cjsaylor/contact-api/controller"
	"github.com/cjsaylor/contact-api/email"
	"github.com/cjsaylor/contact-api/security"
	"github.com/rs/cors"
)

func main() {
	config, err := config.ParseConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	var transport email.Provider
	var validator security.RequestValidator
	if config.TestMode {
		transport = email.NewNullTransport()
		validator = security.NewNullValidator()
		log.Printf("Testmode: ENGAGED")
	} else {
		if config.MailgunPrivateKey != "" {
			client := mailgun.NewMailgun(config.MailgunDomain, config.MailgunPrivateKey, "")
			transport = email.NewMailGunTransport(config.TargetRecipient, client)
		}
		if config.RecaptchaSecretKey != "" {
			client := http.DefaultClient
			validator = security.NewRecaptchaValidator(client, config.RecaptchaSecretKey, config.TestMode)
		}
	}
	mux := http.NewServeMux()
	mux.Handle("/contact", controller.NewContactHTTPHandler(transport, validator))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{config.CORSDomain},
		Debug:          config.TestMode,
	})
	handler := c.Handler(mux)
	log.Printf("Listening on port %v\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), handler))
}
