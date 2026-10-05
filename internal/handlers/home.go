package handlers

import (
	"forum/internal/database"
	"net/http"
)

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !requireGET(w, r) {
		return
	}
	user, err := h.currentUser(r)
	if err != nil {
		serverError(w, err)
		return
	}
	posts, err := database.GetAllPosts(h.DB)
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderPage(w, r, "home", PageData{User: user}, posts)
}
