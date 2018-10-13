package email_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cjsaylor/contact-api/email"
	"github.com/mailgun/mailgun-go"
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

func TestSend(t *testing.T) {
	called := false
	client := mailgun.NewMailgun("testdomain.com", "abc123", "")
	client.SetClient(NewTestClient(func(r *http.Request) *http.Response {
		expectations := map[string]string{
			"from":    "sender@test.com",
			"subject": "some subject",
			"text":    "some body",
			"to":      "test@test.com",
		}
		for key, val := range expectations {
			if r.PostFormValue(key) != val {
				t.Errorf("expected %v got %v", r.PostFormValue(key), val)
			}
		}
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"success":"true"}`)),
			Header:     http.Header{},
		}
	}))
	transport := email.NewMailGunTransport("test@test.com", client)
	err := transport.SendEmail("sender@test.com", "some subject", "some body")
	if err != nil {
		t.Error(err)
	}
	if !called {
		t.Error("expected client send to be triggered")
	}
}
