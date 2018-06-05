package templates

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.New("").ParseFiles("./templates/header.html",
	"./templates/login.html",
	"./templates/footer.html"))

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
