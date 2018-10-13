// Package config is used to configure the application
package config

import (
	"github.com/caarlos0/env"
)

// Configuration holds all application configuration
type Configuration struct {
	Port               int    `env:"PORT" envDefault:"8080"`
	CORSDomain         string `env:"CE_CORS_DOMAIN" envDefault:"http://localhost:1313"`
	TargetRecipient    string `env:"CE_TARGET_EMAIL"`
	RecaptchaSecretKey string `env:"RECAPTCHA_SECRET_KEY"`
	MailgunDomain      string `env:"MG_DOMAIN"`
	MailgunPrivateKey  string `env:"MG_API_KEY"`
	TestMode           bool   `env:"CE_TEST_MODE" envDefault:"true"`
}

// ParseConfiguration retrieves values from environment variables and returns a Configuration struct
func ParseConfiguration() (Configuration, error) {
	config := Configuration{}
	if err := env.Parse(&config); err != nil {
		return config, err
	}
	return config, nil
}
