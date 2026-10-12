package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/utils"
)

func TestProfile(t *testing.T) {
	// Les templates utilisent des chemins relatifs à la racine du projet.
	wd, _ := os.Getwd()
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	db, err := database.Open(filepath.Join(t.TempDir(), "profile.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	hash, _ := utils.HashPassword("old-password")
	id, err := database.CreateUser(db, models.User{Username: "original", Email: "one@example.fr", PasswordHash: hash})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.CreateUser(db, models.User{Username: "occupied", Email: "two@example.fr", PasswordHash: hash}); err != nil {
		t.Fatal(err)
	}
	session := middleware.NewSession(id)
	if err := database.CreateSession(db, session); err != nil {
		t.Fatal(err)
	}
	h := &Handler{DB: db}
	request := func(method string, values url.Values, authenticated bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "/profile", strings.NewReader(values.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if authenticated {
			r.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: session.ID})
		}
		w := httptest.NewRecorder()
		h.ProfileHandler(w, r)
		return w
	}
	if w := request("GET", nil, false); w.Code != 303 || w.Header().Get("Location") != "/login" {
		t.Fatal("anonymous access", w.Code)
	}
	if w := request("GET", nil, true); w.Code != 200 || !strings.Contains(w.Body.String(), "original") || !strings.Contains(w.Body.String(), "csrf_token") {
		t.Fatal("profile rendering", w.Code, w.Body.String())
	}
	values := url.Values{"action": {"username"}, "username": {"renamed"}}
	if w := request("POST", values, true); w.Code != 403 {
		t.Fatal("missing CSRF", w.Code)
	}
	values.Set("csrf_token", profileCSRFToken(session.ID))
	values.Set("username", "occupied")
	if w := request("POST", values, true); !strings.Contains(w.Body.String(), "déjà utilisé") {
		t.Fatal("duplicate username")
	}
	values.Set("username", "renamed")
	if w := request("POST", values, true); w.Code != 303 {
		t.Fatal("rename", w.Code, w.Body.String())
	}
	user, _ := database.GetUserByID(db, int(id))
	if user.Username != "renamed" {
		t.Fatal("username not saved")
	}
	if w := request("GET", nil, true); !strings.Contains(w.Body.String(), "Options du compte de renamed") {
		t.Fatal("greeting not updated")
	}
	values = url.Values{"action": {"password"}, "csrf_token": {profileCSRFToken(session.ID)}, "current_password": {"wrong"}, "new_password": {"new-password"}, "password_confirmation": {"new-password"}}
	if w := request("POST", values, true); !strings.Contains(w.Body.String(), "actuel est incorrect") {
		t.Fatal("incorrect password accepted")
	}
	values.Set("current_password", "old-password")
	values.Set("password_confirmation", "different")
	if w := request("POST", values, true); !strings.Contains(w.Body.String(), "confirmation ne correspond") {
		t.Fatal("mismatch accepted")
	}
	values.Set("password_confirmation", "new-password")
	w := request("POST", values, true)
	if w.Code != 303 || len(w.Result().Cookies()) == 0 {
		t.Fatal("password update", w.Code, w.Body.String())
	}
	user, _ = database.GetUserByID(db, int(id))
	if utils.CheckPassword("new-password", user.PasswordHash) != nil {
		t.Fatal("password not saved")
	}
	if _, err := database.GetSession(db, session.ID); err == nil {
		t.Fatal("old session remains")
	}
	if _, err := database.GetSession(db, w.Result().Cookies()[0].Value); err != nil {
		t.Fatal("new session missing", err)
	}
	if w := request("DELETE", nil, true); w.Code != 405 {
		t.Fatal("unsupported method")
	}
}
