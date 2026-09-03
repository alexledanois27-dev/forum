package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/database"
)

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

	// Charge les templates.
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/home.html",
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Envoie les posts au template.
	if err := tmpl.ExecuteTemplate(w, "layout", posts); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
