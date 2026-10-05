package handlers

import (
	"database/sql"
	"errors"
	"forum/internal/database"
	"net/http"
	"strings"
)

func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	slug := strings.TrimSpace(r.URL.Query().Get("slug"))
	if slug == "" {
		http.Error(w, "Catégorie manquante.", http.StatusBadRequest)
		return
	}
	category, err := database.GetCategoryBySlug(h.DB, slug)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	user, err := h.currentUser(r)
	if err != nil {
		serverError(w, err)
		return
	}
	posts, err := database.GetPostsByCategorySlug(h.DB, slug)
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderPage(w, r, "category", PageData{User: user, Category: category}, posts)
}
