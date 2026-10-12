package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/utils"

	"modernc.org/sqlite"
)

// Le jeton est lié à la session sans exposer son identifiant dans le HTML.
func profileCSRFToken(sessionID string) string {
	sum := sha256.Sum256([]byte("profile-csrf:" + sessionID))
	return hex.EncodeToString(sum[:])
}

func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	user, err := h.currentUser(r)
	if err != nil {
		serverError(w, err)
		return
	}
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionID, err := middleware.GetSessionID(r)
	if err != nil {
		serverError(w, err)
		return
	}
	data := PageData{User: user, CSRFToken: profileCSRFToken(sessionID)}
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("updated") {
		case "username":
			data.ProfileSuccess = "Votre nom d’utilisateur a été mis à jour."
		case "password":
			data.ProfileSuccess = "Votre mot de passe a été modifié."
		}
		h.renderPage(w, r, "profile", data, nil)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide.", http.StatusBadRequest)
		return
	}
	if subtle.ConstantTimeCompare([]byte(r.PostForm.Get("csrf_token")), []byte(data.CSRFToken)) != 1 {
		http.Error(w, "Formulaire expiré ou invalide. Rechargez la page.", http.StatusForbidden)
		return
	}
	showError := func(message string) {
		data.ProfileError = message
		h.renderPage(w, r, "profile", data, nil)
	}
	switch r.PostForm.Get("action") {
	case "username":
		username := r.PostForm.Get("username")
		if err := utils.ValidateUsername(username); err != nil {
			showError("Choisissez un nom de 3 à 30 caractères : lettres sans accents, chiffres, tirets et tirets bas.")
			return
		}
		existing, err := database.GetUserByUsername(h.DB, username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			serverError(w, err)
			return
		}
		if existing != nil && existing.ID != user.ID {
			showError("Ce nom d’utilisateur est déjà utilisé.")
			return
		}
		if err := database.UpdateUsername(h.DB, int(user.ID), username); err != nil {
			var sqliteErr *sqlite.Error
			if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
				showError("Ce nom d’utilisateur est déjà utilisé.")
			} else {
				serverError(w, err)
			}
			return
		}
		http.Redirect(w, r, "/profile?updated=username", http.StatusSeeOther)
	case "password":
		if err := utils.CheckPassword(r.PostForm.Get("current_password"), user.PasswordHash); err != nil {
			showError("Le mot de passe actuel est incorrect.")
			return
		}
		password := r.PostForm.Get("new_password")
		if err := utils.ValidatePassword(password); err != nil {
			showError("Le nouveau mot de passe doit contenir au moins 8 caractères et ne pas dépasser 72 octets.")
			return
		}
		if password != r.PostForm.Get("password_confirmation") {
			showError("La confirmation ne correspond pas au nouveau mot de passe.")
			return
		}
		hash, err := utils.HashPassword(password)
		if err != nil {
			serverError(w, err)
			return
		}
		session := middleware.NewSession(user.ID)
		// Modification et révocation atomiques : aucune ancienne session ne subsiste.
		if err := database.UpdatePasswordAndSession(h.DB, int(user.ID), user.PasswordHash, hash, session); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				showError("Votre compte a changé entre-temps. Rechargez la page.")
			} else {
				serverError(w, err)
			}
			return
		}
		middleware.SetSessionCookie(w, session)
		http.Redirect(w, r, "/profile?updated=password", http.StatusSeeOther)
	default:
		http.Error(w, "Action inconnue.", http.StatusBadRequest)
	}
}
