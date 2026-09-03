package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/utils"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/login.html",
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

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validation de l'email
		if err := utils.ValidateEmail(email); err != nil {
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		// Recherche de l'utilisateur
		user, err := database.GetUserByEmail(h.DB, email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Vérification du mot de passe
		if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Création de la session et du cookie
		_, err = middleware.CreateSession(h.DB, w, int64(user.ID))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
