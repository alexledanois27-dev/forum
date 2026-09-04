package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
)

type PostView struct {
	models.Post

	Author       string
	LikeCount    int
	DislikeCount int
}

type HomePageData struct {
	User  *models.User
	Posts []PostView
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifie l'URL
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Accepte uniquement GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Utilisateur connecté (optionnel)
	var currentUser *models.User

	if session, err := middleware.GetCurrentSession(h.DB, r); err == nil {
		if user, err := database.GetUserByID(h.DB, int(session.UserID)); err == nil {
			currentUser = user
		}
	}

	// Récupère tous les posts
	posts, err := database.GetAllPosts(h.DB)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var postViews []PostView

	for _, post := range posts {

		view := PostView{
			Post: post,
		}

		// Auteur
		if author, err := database.GetUserByID(h.DB, int(post.UserID)); err == nil {
			view.Author = author.Username
		}

		// Nombre de likes
		if likes, err := database.CountPostLikes(h.DB, int(post.ID)); err == nil {
			view.LikeCount = likes
		}

		// Nombre de dislikes
		if dislikes, err := database.CountPostDislikes(h.DB, int(post.ID)); err == nil {
			view.DislikeCount = dislikes
		}

		postViews = append(postViews, view)
	}

	data := HomePageData{
		User:  currentUser,
		Posts: postViews,
	}

	// Charge les templates
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/home.html",
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Affiche la page
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
