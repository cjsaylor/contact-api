package security

import (
	"log"
	"net/http"
)

type nullValidator struct{}

func (n *nullValidator) ValidateRequest(r *http.Request) error {
	log.Println("Request not validated.")
	return nil
}

func NewNullValidator() *nullValidator {
	return &nullValidator{}
}
