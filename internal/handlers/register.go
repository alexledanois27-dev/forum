package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

		if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validation
		if err := utils.ValidateUsername(username); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.ValidateEmail(email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.ValidatePassword(password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash du mot de passe
		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}

		_, err = database.CreateUser(h.DB, user)
		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusBadRequest)
			return
		}

		// L'utilisateur devra ensuite se connecter.
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
