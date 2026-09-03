package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
)

func (h *Handler) CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Vérifie que l'utilisateur est connecté.
	session, err := middleware.GetCurrentSession(h.DB, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	comment := models.Comment{
		PostID:  postID,
		UserID:  session.UserID,
		Content: content,
	}

	_, err = database.CreateComment(h.DB, comment)
	if err != nil {
		http.Error(w, "Unable to create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
