package handlers

import (
	"database/sql"
	"errors"
	"forum/internal/database"
	"net/http"
	"strconv"
)

func (h *Handler) PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Identifiant du post invalide.", http.StatusBadRequest)
		return
	}
	post, err := database.GetPostByID(h.DB, id)
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
	view, err := h.postView(*post, user)
	if err != nil {
		serverError(w, err)
		return
	}
	comments, err := database.GetCommentsByPost(h.DB, id)
	if err != nil {
		serverError(w, err)
		return
	}
	data := PageData{User: user, Post: view, BackURL: safeReturnTo(r.URL.Query().Get("return_to"), "/")}
	for _, comment := range comments {
		view, err := h.commentView(comment, user)
		if err != nil {
			serverError(w, err)
			return
		}
		data.Comments = append(data.Comments, view)
	}
	h.renderPage(w, r, "post-detail", data, nil)
}
