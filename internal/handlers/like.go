package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/middleware"
)

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
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

	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil || (value != 1 && value != -1) {
		http.Error(w, "Invalid reaction value", http.StatusBadRequest)
		return
	}

	// Réaction sur un post
	if postID := r.FormValue("post_id"); postID != "" {
		id, err := strconv.Atoi(postID)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		if err := database.SetPostReaction(h.DB, int(session.UserID), id, value); err != nil {
			http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
			return
		}
	}

	// Réaction sur un commentaire
	if commentID := r.FormValue("comment_id"); commentID != "" {
		id, err := strconv.Atoi(commentID)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		if err := database.SetCommentReaction(h.DB, int(session.UserID), id, value); err != nil {
			http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
