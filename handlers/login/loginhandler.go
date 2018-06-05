package login

import (
	"net/http"

	"github.com/a-h/gauthmiddleware/logger"
	"github.com/a-h/gauthmiddleware/session"
	"github.com/a-h/gauthmiddleware/tokenverifier"
)

const pkg = "github.com/a-h/gauthmiddleware/handlers/login"

// Handler renders the logon screen if you're not logged on, or passes you through to the
// expected content.
type Handler struct {
	Session       session.Session
	TokenVerifier tokenverifier.TokenVerifier
	RenderLogin   http.HandlerFunc
	Next          http.Handler
}

// NewHandler creates an instance of the LoginHandler middleware.
func NewHandler(session session.Session,
	tokenVerifier tokenverifier.TokenVerifier,
	loginRenderer http.HandlerFunc,
	next http.Handler) *Handler {
	return &Handler{
		Session:       session,
		TokenVerifier: tokenVerifier,
		RenderLogin:   loginRenderer,
		Next:          next,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Retrieve the token from Google and validate it against our requirements.
		r.ParseForm()
		idToken := r.FormValue("id_token")

		_, err := h.TokenVerifier.ValidateToken(idToken)
		if err != nil {
			logger.For(pkg, "ServeHTTP").WithField("idToken", idToken).WithError(err).Error("Invalid token")
			http.Error(w, "The presented claim is invalid.", http.StatusInternalServerError)
			return
		}
	}
	isValid, email, err := h.Session.Validate(r)
	if err != nil {
		logger.For(pkg, "ServeHTTP").WithField("email", email).WithError(err).Error("Invalid session")
		http.Error(w, "Unable to validate session.", http.StatusInternalServerError)
		return
	}
	if !isValid {
		h.RenderLogin(w, r)
		return
	}
	logger.For(pkg, "ServeHTTP").WithField("email", email).WithField("url", r.URL.Path).Info("Accessing")
	h.Next.ServeHTTP(w, r)
}
