package controller

import (
	"fmt"
	"net/http"

	"github.com/cjsaylor/contact-api/security"

	"github.com/cjsaylor/contact-api/email"
)

type contactHandler struct {
	emailprovider email.Provider
	validator     security.RequestValidator
}

func NewContactHTTPHandler(provider email.Provider, validator security.RequestValidator) *contactHandler {
	return &contactHandler{
		emailprovider: provider,
		validator:     validator,
	}
}

type input struct {
	sender  string
	subject string
	body    string
}

func inputFromRequest(r *http.Request) (input, error) {
	in := input{}
	if err := r.ParseForm(); err != nil {
		return in, err
	}
	input := map[string]string{
		"sender":  "",
		"subject": "",
		"body":    "",
	}
	required := map[string]interface{}{
		"sender":  nil,
		"subject": nil,
		"body":    nil,
	}
	for key, value := range r.PostForm {
		if _, ok := required[key]; ok {
			input[key] = value[0]
		}
	}
	for key, value := range input {
		if value == "" {
			return in, fmt.Errorf("missing required parameter: %v", key)
		}
	}
	in.sender = input["sender"]
	in.subject = input["subject"]
	in.body = input["body"]
	return in, nil
}

func (c contactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	in, err := inputFromRequest(r)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	if err := c.validator.ValidateRequest(r); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := c.emailprovider.SendEmail(in.sender, in.subject, in.body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to send message."))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
