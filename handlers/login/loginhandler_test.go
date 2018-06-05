package login

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/a-h/gauthmiddleware/session"
	"github.com/a-h/gauthmiddleware/tokenverifier"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name                  string
		request               http.Request
		session               session.Session
		tokenVerifier         func(idToken string) (claim *tokenverifier.Claim, err error)
		expectedNextCalled    bool
		expectedContent       string
		expectedLoginRendered bool
	}{
		{
			name:    "not having a valid session shows the login screen",
			session: mockSession{},
			request: http.Request{
				URL:    &url.URL{Path: "/"},
				Method: "GET",
			},
			expectedNextCalled:    false,
			expectedContent:       "You must login",
			expectedLoginRendered: true,
		},
		{
			name: "having an invalid session shows the login screen",
			session: mockSession{
				validateResponse:             false,
				validateEmailAddressResponse: "marr@example.com",
			},
			request: http.Request{
				URL:    &url.URL{Path: "/"},
				Method: "GET",
			},
			expectedNextCalled:    false,
			expectedContent:       "You must login",
			expectedLoginRendered: true,
		},
		{
			name: "having a valid session from an allowed domain shows the login screen",
			session: mockSession{
				validateResponse:             true,
				validateEmailAddressResponse: "marr@example.com",
			},
			request: http.Request{
				URL:    &url.URL{Path: "/"},
				Method: "GET",
			},
			expectedNextCalled:    true,
			expectedContent:       "Actual content",
			expectedLoginRendered: false,
		},
		{
			name: "POSTing validates the provided auth token - incorrect email address domain",
			tokenVerifier: func(idToken string) (claim *tokenverifier.Claim, err error) {
				if idToken != "the_id_token" {
					t.Errorf("POSTing - expected the_id_token")
				}
				return nil, errors.New("the user is not on the correct GSuite domain")
			},
			session: mockSession{
				validateResponse:             true,
				validateEmailAddressResponse: "marr@example.net",
			},
			request: http.Request{
				URL:    &url.URL{Path: "/"},
				Method: "POST",
				Form: url.Values{
					"id_token": []string{"the_id_token"},
				},
			},
			expectedNextCalled:    false,
			expectedContent:       "The presented claim is invalid.",
			expectedLoginRendered: false,
		},
		{
			name: "POSTing validates the provided auth token - valid email address domain",
			tokenVerifier: func(idToken string) (claim *tokenverifier.Claim, err error) {
				if idToken != "the_id_token" {
					t.Errorf("POSTing - expected the_id_token")
				}
				return &tokenverifier.Claim{
					Email: "marr@example.com",
				}, nil
			},
			session: mockSession{
				validateResponse:             true,
				validateEmailAddressResponse: "marr@example.com",
			},
			request: http.Request{
				URL:    &url.URL{Path: "/"},
				Method: "POST",
				Form: url.Values{
					"id_token": []string{"the_id_token"},
				},
			},
			expectedNextCalled:    true,
			expectedContent:       "Actual content",
			expectedLoginRendered: false,
		},
	}

	for _, test := range tests {
		var actualNextCalled bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualNextCalled = true
			w.Write([]byte("Actual content"))
		})
		var actualLoginRendered bool
		loginRenderer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualLoginRendered = true
			w.Write([]byte("You must login"))
		})
		mtv := mockTokenVerifier{validator: test.tokenVerifier}
		h := NewHandler(test.session, mtv, loginRenderer, next)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, &test.request)

		if test.expectedLoginRendered != actualLoginRendered {
			t.Errorf("%s: expected login rendered to be %v, but was %v", test.name, test.expectedLoginRendered, actualLoginRendered)
		}
		if test.expectedNextCalled != actualNextCalled {
			t.Errorf("%s: expected next called of %v, but was %v", test.name, test.expectedNextCalled, actualNextCalled)
		}

		actualBody, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			t.Fatalf("%s: failed to read body of result: %v", test.name, err)
		}
		if !bytes.Contains(actualBody, []byte(test.expectedContent)) {
			t.Errorf("%s: expected body to contain %s, but it didn't: %s", test.name, test.expectedContent, string(actualBody))
		}
	}
}

type mockTokenVerifier struct {
	validator func(idToken string) (claim *tokenverifier.Claim, err error)
}

func (m mockTokenVerifier) ValidateToken(idToken string) (claim *tokenverifier.Claim, err error) {
	return m.validator(idToken)
}
