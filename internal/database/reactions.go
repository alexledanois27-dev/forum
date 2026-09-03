package database

import "database/sql"

// SetPostReaction creates or updates a user's reaction to a post.
func SetPostReaction(db *sql.DB, userID, postID, value int) error {
	if value != 1 && value != -1 {
		return ErrInvalidReactionValue
	}

	_, err := db.Exec(`
		INSERT INTO post_reactions (user_id, post_id, value)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id)
		DO UPDATE SET value = excluded.value, created_at = CURRENT_TIMESTAMP
	`, userID, postID, value)

	return err
}

// RemovePostReaction removes a user's reaction from a post.
func RemovePostReaction(db *sql.DB, userID, postID int) error {
	_, err := db.Exec(`
		DELETE FROM post_reactions
		WHERE user_id = ? AND post_id = ?
	`, userID, postID)

	return err
}

// GetPostReaction returns a user's reaction to a post.
func GetPostReaction(db *sql.DB, userID, postID int) (int, error) {
	var value int

	err := db.QueryRow(`
		SELECT value
		FROM post_reactions
		WHERE user_id = ? AND post_id = ?
	`, userID, postID).Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// CountPostLikes returns the number of likes on a post.
func CountPostLikes(db *sql.DB, postID int) (int, error) {
	return countPostReactions(db, postID, 1)
}

// CountPostDislikes returns the number of dislikes on a post.
func CountPostDislikes(db *sql.DB, postID int) (int, error) {
	return countPostReactions(db, postID, -1)
}

// countPostReactions counts post reactions for a given value.
func countPostReactions(db *sql.DB, postID, value int) (int, error) {
	var count int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM post_reactions
		WHERE post_id = ? AND value = ?
	`, postID, value).Scan(&count)

	return count, err
}

// SetCommentReaction creates or updates a user's reaction to a comment.
func SetCommentReaction(db *sql.DB, userID, commentID, value int) error {
	if value != 1 && value != -1 {
		return ErrInvalidReactionValue
	}

	_, err := db.Exec(`
		INSERT INTO comment_reactions (user_id, comment_id, value)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, comment_id)
		DO UPDATE SET value = excluded.value, created_at = CURRENT_TIMESTAMP
	`, userID, commentID, value)

	return err
}

// RemoveCommentReaction removes a user's reaction from a comment.
func RemoveCommentReaction(db *sql.DB, userID, commentID int) error {
	_, err := db.Exec(`
		DELETE FROM comment_reactions
		WHERE user_id = ? AND comment_id = ?
	`, userID, commentID)

	return err
}

// GetCommentReaction returns a user's reaction to a comment.
func GetCommentReaction(db *sql.DB, userID, commentID int) (int, error) {
	var value int

	err := db.QueryRow(`
		SELECT value
		FROM comment_reactions
		WHERE user_id = ? AND comment_id = ?
	`, userID, commentID).Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// CountCommentLikes returns the number of likes on a comment.
func CountCommentLikes(db *sql.DB, commentID int) (int, error) {
	return countCommentReactions(db, commentID, 1)
}

// CountCommentDislikes returns the number of dislikes on a comment.
func CountCommentDislikes(db *sql.DB, commentID int) (int, error) {
	return countCommentReactions(db, commentID, -1)
}

// countCommentReactions counts comment reactions for a given value.
func countCommentReactions(db *sql.DB, commentID, value int) (int, error) {
	var count int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM comment_reactions
		WHERE comment_id = ? AND value = ?
	`, commentID, value).Scan(&count)

	return count, err
}
