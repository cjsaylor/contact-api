package security

import (
	"net/http"
)

type RequestValidator interface {
	ValidateRequest(r *http.Request) error
}
