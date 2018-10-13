package security

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

type recaptchaValidator struct {
	client    *http.Client
	secretKey string
	debug     bool
}

type response struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Errors      []string  `json:"error-codes"`
}

func (r *recaptchaValidator) ValidateRequest(req *http.Request) error {
	token := req.PostFormValue("g-recaptcha-response")
	params := url.Values{
		"secret":   {r.secretKey},
		"response": {token},
	}
	resp, err := r.client.PostForm(verifyURL, params)
	if err != nil {
		return err
	}
	result := response{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	if !result.Success {
		if r.debug {
			log.Println(result.Errors)
		}
		return fmt.Errorf("Error validating recaptcha token")
	}
	return nil
}

// NewRecaptchaValidator returns a validator for recaptcha response validation.
func NewRecaptchaValidator(client *http.Client, secretKey string, debugMode bool) *recaptchaValidator {
	return &recaptchaValidator{
		client:    client,
		secretKey: secretKey,
		debug:     debugMode,
	}
}
