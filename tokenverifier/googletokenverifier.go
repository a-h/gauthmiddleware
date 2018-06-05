package tokenverifier

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// A GoogleTokenVerifier verifies Tokens with Google.
type GoogleTokenVerifier struct {
	AllowedDomains []string
}

// ValidateToken retrieves a claim from Google and validates it using Google's rules.
func (verifier GoogleTokenVerifier) ValidateToken(idToken string) (claim *Claim, err error) {
	claim, err = verifier.GetClaim(idToken)
	if err != nil {
		return
	}
	_, err = verifier.IsClaimValid(claim)
	return claim, err
}

// GetClaim returns a Claim from Google, using the id_token presented by the
// Google authentication system.
func (verifier GoogleTokenVerifier) GetClaim(idToken string) (*Claim, error) {
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + url.QueryEscape(idToken)
	body, err := getResponse(url)
	if err != nil {
		return nil, fmt.Errorf("GoogleTokenVerifier: failed to get token with error: %v", err)
	}
	return NewClaim(body)
}

func getResponse(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// IsClaimValid validates a claim by checking that it's not expired, the issuer
// was Google and that the user's email address has been verified.
func (verifier GoogleTokenVerifier) IsClaimValid(claim *Claim) (ok bool, err error) {
	expiry, expiryErr := strconv.Atoi(claim.Expiry)
	emailVerified, emailVerifiedErr := strconv.ParseBool(claim.EmailVerified)

	validation := []struct {
		name               string
		validationFunction func() bool
	}{
		{"email ok", func() bool { return claim.Email != "" }},
		{"email verified ok", func() bool { return emailVerifiedErr == nil && emailVerified }},
		{"expiry is number", func() bool { return expiryErr == nil }},
		{"expiry ok", func() bool { return time.Unix(int64(expiry), 0).After(time.Now()) }},
		{"issuer ok", func() bool {
			return claim.Issuer == "https://accounts.google.com" || claim.Issuer == "accounts.google.com"
		}},
		{"domain ok", func() bool {
			if len(verifier.AllowedDomains) == 0 {
				return true
			}
			for _, ad := range verifier.AllowedDomains {
				if claim.HD == ad {
					return true
				}
			}
			return false
		}},
	}

	var errorMessage bytes.Buffer
	ok = true
	for _, v := range validation {
		valid := v.validationFunction()
		errorMessage.WriteString(v.name + " " + strconv.FormatBool(valid) + "\n")
		ok = ok && valid
	}
	if !ok {
		err = errors.New(errorMessage.String())
	}
	return
}
