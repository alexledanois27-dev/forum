package database

import (
	"database/sql"

	"forum/internal/models"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, user models.User) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO users (username, email, password_hash)
		VALUES (?, ?, ?)
	`, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUserByID returns a user by its ID.
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User

	err := db.QueryRow(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail returns a user by its email address.
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	var user models.User

	err := db.QueryRow(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername returns a user by its username.
func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User

	err := db.QueryRow(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUsername updates a user's username.
func UpdateUsername(db *sql.DB, userID int, username string) error {
	_, err := db.Exec(`
		UPDATE users
		SET username = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, username, userID)

	return err
}

// UpdateUserPassword updates a user's password hash.
func UpdateUserPassword(db *sql.DB, userID int, passwordHash string) error {
	_, err := db.Exec(`
		UPDATE users
		SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, passwordHash, userID)

	return err
}

// DeleteUser removes a user from the database.
func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, userID)

	return err
}

// UpdatePasswordAndSession remplace le mot de passe et toutes les sessions
// dans une seule transaction. Le hash précédent protège des mises à jour concurrentes.
func UpdatePasswordAndSession(db *sql.DB, userID int, oldHash, newHash string, session models.Session) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec("UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND password_hash = ?", newHash, userID, oldHash)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return sql.ErrNoRows
	}
	if _, err := tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", session.ID, userID, session.ExpiresAt); err != nil {
		return err
	}
	return tx.Commit()
}
