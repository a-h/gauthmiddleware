package session

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	tests := []struct {
		name          string
		request       func() (*http.Request, error)
		expectedValid bool
		expectedEmail string
	}{
		{
			name: "no session cookie",
			request: func() (*http.Request, error) {
				return http.NewRequest("GET", "http://example.com", nil)
			},
			expectedValid: false,
		},
		{
			name: "valid session cookie",
			request: func() (*http.Request, error) {
				r, err := http.NewRequest("GET", "http://example.com", nil)
				if err != nil {
					return nil, fmt.Errorf("error setting up request: %v", err)
				}
				w := httptest.NewRecorder()
				s := NewGorillaSession([]byte("random_data"), false, "cookie-name")
				err = s.Start(w, r, "test@example.com")
				return r, err
			},
			expectedValid: true,
			expectedEmail: "test@example.com",
		},
	}

	for _, test := range tests {
		s := NewGorillaSession([]byte("random_data"), false, "cookie-name")

		r, err := test.request()
		if err != nil {
			t.Fatalf("%s: error creating test request: %v", test.name, err)
		}

		actualValid, actualEmail, err := s.Validate(r)
		if err != nil {
			t.Fatalf("%s: unexpected error validating the session: %v", test.name, err)
		}
		if test.expectedValid != actualValid {
			t.Errorf("%s: expected valid %v, got %v", test.name, test.expectedValid, actualValid)
		}
		if test.expectedEmail != actualEmail {
			t.Errorf("%s: expected email %v, got %v", test.name, test.expectedEmail, actualEmail)
		}
	}
}
