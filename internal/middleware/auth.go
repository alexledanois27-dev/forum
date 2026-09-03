package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/models"

	"github.com/google/uuid"
)

const (
	SessionCookieName = "session_id"
	SessionDuration   = 24 * time.Hour
)

// NewSession creates a session ready to be stored in the database.
func NewSession(userID int64) models.Session {
	now := time.Now().UTC().Truncate(time.Second)

	return models.Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		ExpiresAt: now.Add(SessionDuration),
		CreatedAt: now,
	}
}

// CreateSession replaces the user's previous session, stores the new one and
// sends its identifier to the browser.
func CreateSession(db *sql.DB, w http.ResponseWriter, userID int64) (*models.Session, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	session := NewSession(userID)
	if err := database.CreateSession(db, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	SetSessionCookie(w, session)

	return &session, nil
}

// SetSessionCookie sends the session identifier to the browser.
func SetSessionCookie(w http.ResponseWriter, session models.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		MaxAge:   int(time.Until(session.ExpiresAt).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetSessionID returns and validates the session identifier stored in the cookie.
func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}

	if cookie.Value == "" {
		return "", errors.New("empty session cookie")
	}

	if _, err := uuid.Parse(cookie.Value); err != nil {
		return "", errors.New("invalid session cookie")
	}

	return cookie.Value, nil
}

// GetCurrentSession returns the valid database session associated with the
// request cookie.
func GetCurrentSession(db *sql.DB, r *http.Request) (*models.Session, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	sessionID, err := GetSessionID(r)
	if err != nil {
		return nil, err
	}

	session, err := database.GetValidSession(db, sessionID, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("get valid session: %w", err)
	}

	return session, nil
}

// ClearSessionCookie asks the browser to immediately delete its session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// DeleteSession removes the server-side session and always clears the browser
// cookie. Calling it without a valid cookie is safe and has no effect.
func DeleteSession(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	ClearSessionCookie(w)

	sessionID, err := GetSessionID(r)
	if err != nil {
		return nil
	}

	if db == nil {
		return errors.New("database is required")
	}

	if err := database.DeleteSession(db, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
