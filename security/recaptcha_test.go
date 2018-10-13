package security_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cjsaylor/contact-api/security"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
func TestRequestValidation(t *testing.T) {
	client := NewTestClient(func(r *http.Request) *http.Response {
		if r.Method != http.MethodPost {
			t.Errorf("expecting %v got %v", http.MethodPost, r.Method)
		}
		if r.PostFormValue("secret") != "abcd" {
			t.Errorf("expecting secret abcd got %v", r.PostFormValue("secret"))
		}
		if r.PostFormValue("response") != "response-1234" {
			t.Errorf("expecting response for server request (response-1234) got %v", r.PostFormValue("response"))
		}
		if r.URL.String() != "https://www.google.com/recaptcha/api/siteverify" {
			t.Errorf("unexpected verify endpoint: %v", r.URL.String())
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"success":"true"}`)),
			Header:     http.Header{},
		}
	})
	recaptcha := security.NewRecaptchaValidator(client, "abcd", false)
	data := url.Values{
		"g-recaptcha-response": {"response-1234"},
	}
	rBody := strings.NewReader(data.Encode())
	serverReq := httptest.NewRequest(http.MethodPost, "/contact", rBody)
	serverReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recaptcha.ValidateRequest(serverReq)
}
