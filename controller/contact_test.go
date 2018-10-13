package controller_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cjsaylor/contact-api/controller"
	"github.com/cjsaylor/contact-api/email"
	"github.com/cjsaylor/contact-api/security"
)

func TestMissingInput(t *testing.T) {
	contact := controller.NewContactHTTPHandler(email.NewNullTransport(), security.NewNullValidator())
	data := url.Values{
		"sender":  []string{"test@test.com"},
		"subject": []string{"some subject"},
		// intentionally neglecting the body required parameter
	}
	rBody := strings.NewReader(data.Encode())
	r := httptest.NewRequest(http.MethodPost, "/contact", rBody)
	w := httptest.NewRecorder()
	contact.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected code %v got %v", http.StatusBadRequest, resp.StatusCode)
	}
	if string(body)[:26] != "missing required parameter" {
		t.Errorf("expected an error message, got %v", string(body)[:26])
	}
}

func TestSuccessfulInput(t *testing.T) {
	contact := controller.NewContactHTTPHandler(email.NewNullTransport(), security.NewNullValidator())
	data := url.Values{
		"sender":  []string{"test@test.com"},
		"subject": []string{"some subject"},
		"body":    []string{"some comment body"},
	}
	rBody := strings.NewReader(data.Encode())
	r := httptest.NewRequest(http.MethodPost, "/contact", rBody)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	contact.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v got %v", http.StatusOK, resp.StatusCode)
	}
}
