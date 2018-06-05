package configuration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Configuration contains the configuration of the application.
type Configuration struct {
	// SessionEncryptionKey is used to encrypt the user session details. It should be 32 bytes of random data.
	SessionEncryptionKey []byte
	// CookieName is the name of the session cookie.
	CookieName string
	// SetSecureFlag sets whether cookies should be issued with the secure flag set.
	// When the secure flag is set, cookies cannot be transmitted over HTTP.
	// SSL must already be in place before this option is set.
	SetSecureFlag bool
	// GoogleAuthClientID is required to enable authentication.
	GoogleAuthClientID string
	// GoogleAllowedDomains are Google GSuite domains which are permitted to access the content.
	GoogleAllowedDomains []string
}

// FromEnvironment loads the configuration using environment variables.
func FromEnvironment() (c Configuration, err error) {
	var errs []string

	var n int
	n, err = base64.RawStdEncoding.Decode(c.SessionEncryptionKey, []byte(os.Getenv("SESSION_ENCRYPTION_KEY")))
	if err != nil {
		errs = append(errs, fmt.Sprintf("SESSION_ENCRYPTION_KEY: not found, or invalid: %v", err))
	}
	if n != 32 {
		errs = append(errs, fmt.Sprintf("SESSION_ENCRYPTION_KEY: expected 32 bytes when base64 decoded"))
	}

	c.CookieName = os.Getenv("COOKIE_NAME")
	if c.CookieName == "" {
		errs = append(errs, fmt.Sprintf("COOKIE_NAME: not set"))
	}

	c.SetSecureFlag, err = strconv.ParseBool(os.Getenv("SET_SECURE_FLAG"))
	if err != nil {
		errs = append(errs, fmt.Sprintf("SET_SECURE_FLAG: not set or invalid value: '%v'", os.Getenv("SET_SECURE_FLAG")))
	}

	c.GoogleAuthClientID = os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	if c.GoogleAuthClientID == "" {
		errs = append(errs, fmt.Sprintf("GOOGLE_AUTH_CLIENT_ID: not set"))
	}

	gad := os.Getenv("GOOGLE_ALLOWED_DOMAINS")
	if gad != "*" {
		c.GoogleAllowedDomains = strings.Split(gad, ",")
		if len(c.GoogleAllowedDomains) == 0 {
			errs = append(errs, fmt.Sprintf("GOOGLE_ALLOWED_DOMAINS: not set"))
		}
	}

	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, ", "))
	}

	return
}
