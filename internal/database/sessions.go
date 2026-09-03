package database

import (
	"database/sql"
	"time"

	"forum/internal/models"
)

// CreateSession creates a new session and replaces any existing session for the user.
func CreateSession(db *sql.DB, session models.Session) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM sessions WHERE user_id = ?`, session.UserID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)
	`, session.ID, session.UserID, session.ExpiresAt); err != nil {
		return err
	}

	return tx.Commit()
}

// GetSession returns a session by its ID.
func GetSession(db *sql.DB, sessionID string) (*models.Session, error) {
	var session models.Session

	err := db.QueryRow(`
		SELECT id, user_id, expires_at, created_at
		FROM sessions
		WHERE id = ?
	`, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetValidSession returns a session only if it has not expired.
func GetValidSession(db *sql.DB, sessionID string, now time.Time) (*models.Session, error) {
	var session models.Session

	err := db.QueryRow(`
		SELECT id, user_id, expires_at, created_at
		FROM sessions
		WHERE id = ? AND expires_at > ?
	`, sessionID, now).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession removes a session by its ID.
func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`
		DELETE FROM sessions
		WHERE id = ?
	`, sessionID)

	return err
}

// DeleteSessionsByUser removes all sessions for a user.
func DeleteSessionsByUser(db *sql.DB, userID int) error {
	_, err := db.Exec(`
		DELETE FROM sessions
		WHERE user_id = ?
	`, userID)

	return err
}

// DeleteExpiredSessions removes all expired sessions.
func DeleteExpiredSessions(db *sql.DB, now time.Time) error {
	_, err := db.Exec(`
		DELETE FROM sessions
		WHERE expires_at <= ?
	`, now)

	return err
}
