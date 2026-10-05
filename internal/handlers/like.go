package handlers

import (
	"database/sql"
	"errors"
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

	postTarget := r.FormValue("post_id")
	commentTarget := r.FormValue("comment_id")
	if (postTarget == "") == (commentTarget == "") {
		http.Error(w, "Choisis un post ou un commentaire.", http.StatusBadRequest)
		return
	}
	fallback := "/"
	// Réaction sur un post
	if postID := r.FormValue("post_id"); postID != "" {
		id, err := strconv.Atoi(postID)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		_, err = database.GetPostByID(h.DB, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
			} else {
				serverError(w, err)
			}
			return
		}
		fallback = "/post-detail?id=" + strconv.Itoa(id)
		if err := database.SetPostReaction(h.DB, int(session.UserID), id, value); err != nil {
			http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
			return
		}
	}

	// Réaction sur un commentaire
	if commentID := r.FormValue("comment_id"); commentID != "" {
		id, err := strconv.Atoi(commentID)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		target, err := database.GetCommentByID(h.DB, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
			} else {
				serverError(w, err)
			}
			return
		}
		fallback = "/post-detail?id=" + strconv.FormatInt(target.PostID, 10)
		if err := database.SetCommentReaction(h.DB, int(session.UserID), id, value); err != nil {
			http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
			return
		}
	}

	redirectBack(w, r, fallback)
}
