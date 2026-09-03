package handlers

import (
	"html/template"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/register.html",
		)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", nil)

	case http.MethodPost:

		// TODO: r.ParseForm()

		// TODO: username := r.FormValue("username")
		// TODO: email := r.FormValue("email")
		// TODO: password := r.FormValue("password")

		// TODO: validation.Validate...

		// TODO: hash.HashPassword(...)

		// TODO: database.CreateUser(...)

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
