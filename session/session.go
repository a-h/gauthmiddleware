package session

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// Session determines how a user is logged in to the system.
type Session interface {
	// ValidateSession validates a session and returns the email address of the
	// user.
	Validate(r *http.Request) (isValid bool, emailAddress string, err error)
	Start(w http.ResponseWriter, r *http.Request, emailAddress string) error
}

// A GorillaSession uses the Gorilla framework to manage the session.
type GorillaSession struct {
	store      sessions.CookieStore
	CookieName string
}

// NewGorillaSession creates a Session which uses Gorilla.
func NewGorillaSession(encryptionKey []byte, setSecureFlag bool, cookieName string) *GorillaSession {
	store := sessions.NewCookieStore(encryptionKey)
	store.Options = &sessions.Options{
		HttpOnly: true,
		Secure:   setSecureFlag,
	}
	return &GorillaSession{
		store:      *store,
		CookieName: cookieName,
	}
}

// Start starts off a session by adding the emailAddress value to an
// encrypted cookie.
func (gs GorillaSession) Start(w http.ResponseWriter, r *http.Request, emailAdress string) error {
	session, err := gs.store.Get(r, gs.CookieName)
	if err != nil {
		return err
	}
	session.Values["emailAddress"] = emailAdress
	return session.Save(r, w)
}

// Validate checks whether the session is valid. If it isn't, it will
// redirect the user to the logon screen.
func (gs GorillaSession) Validate(r *http.Request) (isValid bool, emailAddress string, err error) {
	session, err := gs.store.Get(r, gs.CookieName)
	if err != nil {
		err = fmt.Errorf("GorillaSession.Validate: failed to get the cookie from the store: %v", err)
		return
	}
	emailAddress, isValid = session.Values["emailAddress"].(string)
	return
}
