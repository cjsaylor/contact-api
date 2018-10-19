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
	if config.MailgunPrivateKey != "" && !config.TestMode {
		client := mailgun.NewMailgun(config.MailgunDomain, config.MailgunPrivateKey, "")
		transport = email.NewMailGunTransport(config.TargetRecipient, client)
	} else {
		transport = email.NewNullTransport()
		log.Println("Null email transport initialized. Not sending emails.")
	}
	if config.RecaptchaSecretKey != "" && !config.TestMode {
		client := http.DefaultClient
		validator = security.NewRecaptchaValidator(client, config.RecaptchaSecretKey, config.TestMode)
	} else {
		validator = security.NewNullValidator()
		log.Println("Null security validator initialized. Not validating tokens.")
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
