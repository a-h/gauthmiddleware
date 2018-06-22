package templates

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func init() {
	templates = template.New("")
	template.Must(templates.New("header.html").Parse(string(MustAsset("templates/header.html"))))
	template.Must(templates.New("login.html").Parse(string(MustAsset("templates/login.html"))))
	template.Must(templates.New("footer.html").Parse(string(MustAsset("templates/footer.html"))))
}

// LoginModel is the data required to render the Login screen.
type LoginModel struct {
	GoogleAuthClientID string
}

// RenderLogin renders the login template.
func RenderLogin(w http.ResponseWriter, model LoginModel) error {
	return Render(w, "login.html", model)
}

// Render template to HTTP.
func Render(w http.ResponseWriter, templateName string, model interface{}) (err error) {
	err = templates.ExecuteTemplate(w, templateName, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
