package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/database"
	"forum/internal/models"
)

type HomePageData struct {
	User  *models.User
	Posts []models.Post
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie que l'utilisateur demande bien la page d'accueil.
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Accepte uniquement les requêtes GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupère tous les posts.
	posts, err := database.GetAllPosts(h.DB)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Récupère la session de l'utilisateur.
	session, err := database.GetCurrentSession(h.DB, r)
	if err != nil {
		session = nil
	}

	var user *models.User

	// Si une session existe, récupère l'utilisateur correspondant.
	if session != nil {
		user, err = database.GetUserByID(h.DB, session.UserID)
		if err != nil {
			user = nil
		}
	}

	// Prépare les données envoyées aux templates.
	data := HomePageData{
		User:  user,
		Posts: posts,
	}

	// Charge les templates.
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/home.html",
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Envoie les données au template.
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
