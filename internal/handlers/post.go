package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
)

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	// L'utilisateur doit être connecté.
	session, err := middleware.GetCurrentSession(h.DB, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {

	case http.MethodGet:
		categories, err := database.GetAllCategories(h.DB)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/post.html",
		)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "layout", categories); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		post := models.Post{
			UserID:  session.UserID,
			Title:   title,
			Content: content,
		}

		var categoryIDs []int

		for _, value := range r.Form["categories"] {
			id, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			categoryIDs = append(categoryIDs, id)
		}

		_, err = database.CreatePost(h.DB, post, categoryIDs)
		if err != nil {
			http.Error(w, "Unable to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
