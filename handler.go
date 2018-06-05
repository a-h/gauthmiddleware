package gauthmiddleware

import (
	"net/http"

	"github.com/a-h/gauthmiddleware/configuration"
	"github.com/a-h/gauthmiddleware/handlers/login"
	"github.com/a-h/gauthmiddleware/session"
	"github.com/a-h/gauthmiddleware/templates"
	"github.com/a-h/gauthmiddleware/tokenverifier"
)

// New starts up the middleware, loading all configuration from environment variables.
func New(next http.Handler) (h http.Handler, err error) {
	conf, err := configuration.FromEnvironment()
	if err != nil {
		return
	}
	h = NewWithConfiguration(conf, next)
	return
}

// NewWithConfiguration starts up the GAuth middleware using the provided configuration.
func NewWithConfiguration(conf configuration.Configuration, next http.Handler) http.Handler {
	session := session.NewGorillaSession(conf.SessionEncryptionKey, conf.SetSecureFlag, conf.CookieName)
	lr := func(w http.ResponseWriter, r *http.Request) {
		templates.RenderLogin(w, templates.LoginModel{
			GoogleAuthClientID: conf.GoogleAuthClientID,
		})
	}
	tv := tokenverifier.GoogleTokenVerifier{
		AllowedDomains: conf.GoogleAllowedDomains,
	}
	return login.NewHandler(session, tv, lr, next)
}
