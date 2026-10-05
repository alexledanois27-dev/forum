package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" || utf8.RuneCountInString(content) > 800 {
		http.Error(w, "Le commentaire doit contenir entre 1 et 800 caractères.", http.StatusBadRequest)
		return
	}
	if _, err := database.GetPostByID(h.DB, int(postID)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			serverError(w, err)
		}
		return
	}

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

	redirectBack(w, r, "/post-detail?id="+strconv.FormatInt(postID, 10))
}
