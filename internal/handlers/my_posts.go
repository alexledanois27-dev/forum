package handlers

import (
	"forum/internal/database"
	"net/http"
)

func (h *Handler) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	user, err := h.currentUser(r)
	if err != nil {
		serverError(w, err)
		return
	}
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	posts, err := database.GetPostsByUser(h.DB, int(user.ID))
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderPage(w, r, "my-posts", PageData{User: user}, posts)
}
