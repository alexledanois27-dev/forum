package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	_ "time/tzdata"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
)

type PostView struct {
	models.Post
	Author                                              string
	LikeCount, DislikeCount, CommentCount, UserReaction int
}

type CommentView struct {
	models.Comment
	Author                                string
	LikeCount, DislikeCount, UserReaction int
}

type PageData struct {
	User       *models.User
	Posts      []PostView
	Post       PostView
	Comments   []CommentView
	Category   *models.Category
	Categories []models.Category
	ReturnTo   string
	BackURL    string
}

// Conserve les dates de la base ; convertit uniquement les copies affichées.
var displayLocation = func() *time.Location {
	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	return location
}()

func (h *Handler) currentUser(r *http.Request) (*models.User, error) {
	if _, err := r.Cookie(middleware.SessionCookieName); errors.Is(err, http.ErrNoCookie) {
		return nil, nil
	}
	// GetSessionID ne consulte pas la base : un cookie mal formé est anonyme.
	if _, err := middleware.GetSessionID(r); err != nil {
		return nil, nil
	}
	session, err := middleware.GetCurrentSession(h.DB, r)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user, err := database.GetUserByID(h.DB, int(session.UserID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (h *Handler) postView(post models.Post, user *models.User) (PostView, error) {
	view := PostView{Post: post}
	view.CreatedAt = post.CreatedAt.In(displayLocation)
	view.UpdatedAt = post.UpdatedAt.In(displayLocation)
	author, err := database.GetUserByID(h.DB, int(post.UserID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return view, err
	}
	if author != nil {
		view.Author = author.Username
	}
	if view.LikeCount, err = database.CountPostLikes(h.DB, int(post.ID)); err != nil {
		return view, err
	}
	if view.DislikeCount, err = database.CountPostDislikes(h.DB, int(post.ID)); err != nil {
		return view, err
	}
	if view.CommentCount, err = database.CountPostComments(h.DB, int(post.ID)); err != nil {
		return view, err
	}
	if user != nil {
		view.UserReaction, err = database.GetPostReaction(h.DB, int(user.ID), int(post.ID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return view, err
		}
	}
	return view, nil
}

func (h *Handler) commentView(comment models.Comment, user *models.User) (CommentView, error) {
	view := CommentView{Comment: comment}
	view.CreatedAt = comment.CreatedAt.In(displayLocation)
	view.UpdatedAt = comment.UpdatedAt.In(displayLocation)
	author, err := database.GetUserByID(h.DB, int(comment.UserID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return view, err
	}
	if author != nil {
		view.Author = author.Username
	}
	if view.LikeCount, err = database.CountCommentLikes(h.DB, int(comment.ID)); err != nil {
		return view, err
	}
	if view.DislikeCount, err = database.CountCommentDislikes(h.DB, int(comment.ID)); err != nil {
		return view, err
	}
	if user != nil {
		view.UserReaction, err = database.GetCommentReaction(h.DB, int(user.ID), int(comment.ID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return view, err
		}
	}
	return view, nil
}

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, page string, data PageData, posts []models.Post) {
	var err error
	data.Categories, err = database.GetAllCategories(h.DB)
	if err != nil {
		serverError(w, err)
		return
	}
	data.ReturnTo = safeReturnTo(r.URL.RequestURI(), "/")
	for _, post := range posts {
		view, err := h.postView(post, data.User)
		if err != nil {
			serverError(w, err)
			return
		}
		data.Posts = append(data.Posts, view)
	}
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/"+page+".html")
	if err != nil {
		serverError(w, err)
		return
	}
	// Rend en mémoire pour ne pas envoyer une page incomplète en cas d'erreur.
	var buffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buffer, "layout", data); err != nil {
		serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buffer.Bytes())
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("handler: %v", err)
	http.Error(w, "Une erreur interne est survenue.", http.StatusInternalServerError)
}

func requireGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		return true
	}
	w.Header().Set("Allow", http.MethodGet)
	http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
	return false
}

// Autorise seulement les pages de lecture connues, jamais une URL externe.
func safeReturnTo(raw, fallback string) string {
	if strings.ContainsAny(raw, "\\\r\n") {
		return fallback
	}
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || u.Host != "" || u.Opaque != "" {
		return fallback
	}
	switch u.Path {
	case "/", "/category", "/my-posts", "/post-detail":
		return u.String()
	default:
		return fallback
	}
}

func redirectBack(w http.ResponseWriter, r *http.Request, fallback string) {
	target := r.PostForm.Get("return_to")
	if target == "" {
		target = r.Referer()
	}
	// Les Referer absolus sont ignorés ; les formulaires transmettent return_to.
	http.Redirect(w, r, safeReturnTo(target, fallback), http.StatusSeeOther)
}
