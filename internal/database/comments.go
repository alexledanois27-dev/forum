package database

import (
	"database/sql"

	"forum/internal/models"
)

// CreateComment inserts a new comment into the database.
func CreateComment(db *sql.DB, comment models.Comment) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)
	`, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetCommentByID returns a comment by its ID.
func GetCommentByID(db *sql.DB, commentID int) (*models.Comment, error) {
	var comment models.Comment

	err := db.QueryRow(`
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE id = ?
	`, commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetCommentsByPost returns all comments for a post.
func GetCommentsByPost(db *sql.DB, postID int) ([]models.Comment, error) {
	rows, err := db.Query(`
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// UpdateComment updates the content of an existing comment.
func UpdateComment(db *sql.DB, commentID int, content string) error {
	_, err := db.Exec(`
		UPDATE comments
		SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, content, commentID)

	return err
}

// DeleteComment removes a comment from the database.
func DeleteComment(db *sql.DB, commentID int) error {
	_, err := db.Exec(`
		DELETE FROM comments
		WHERE id = ?
	`, commentID)

	return err
}
