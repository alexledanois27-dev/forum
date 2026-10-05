package handlers

import (
	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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

		user, err := database.GetUserByID(h.DB, int(session.UserID))
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

		data := struct {
			User       *models.User
			Categories []models.Category
		}{
			User:       user,
			Categories: categories,
		}

		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		if title == "" || utf8.RuneCountInString(title) > 150 {
			http.Error(w, "Le titre doit contenir entre 1 et 150 caractères.", http.StatusBadRequest)
			return
		}

		if content == "" || utf8.RuneCountInString(content) > 10000 {
			http.Error(w, "Le post doit contenir entre 1 et 10 000 caractères.", http.StatusBadRequest)
			return
		}

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
