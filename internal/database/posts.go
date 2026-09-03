package database

import (
	"database/sql"

	"forum/internal/models"
)

// CreatePost creates a new post and links it to the selected categories.
func CreatePost(db *sql.DB, post models.Post, categoryIDs []int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO posts (user_id, title, content)
		VALUES (?, ?, ?)
	`, post.UserID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, categoryID := range categoryIDs {
		if _, err := tx.Exec(`
			INSERT INTO post_categories (post_id, category_id)
			VALUES (?, ?)
		`, postID, categoryID); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return postID, nil
}

// GetAllPosts returns all posts ordered from newest to oldest.
func GetAllPosts(db *sql.DB) ([]models.Post, error) {
	return queryPosts(db, `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
}

// GetPostByID returns a post by its ID.
func GetPostByID(db *sql.DB, postID int) (*models.Post, error) {
	var post models.Post

	err := db.QueryRow(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		WHERE id = ?
	`, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetPostsByUser returns all posts created by a user.
func GetPostsByUser(db *sql.DB, userID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetPostsByCategory returns all posts linked to a category.
func GetPostsByCategory(db *sql.DB, categoryID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at
		FROM posts p
		INNER JOIN post_categories pc ON pc.post_id = p.id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetPostsByCategorySlug returns all posts linked to a category slug.
func GetPostsByCategorySlug(db *sql.DB, slug string) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at
		FROM posts p
		INNER JOIN post_categories pc ON pc.post_id = p.id
		INNER JOIN categories c ON c.id = pc.category_id
		WHERE c.slug = ?
		ORDER BY p.created_at DESC
	`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetLikedPostsByUser returns all posts liked by a user.
func GetLikedPostsByUser(db *sql.DB, userID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at
		FROM posts p
		INNER JOIN post_reactions pr ON pr.post_id = p.id
		WHERE pr.user_id = ? AND pr.value = 1
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// UpdatePost updates the title and content of a post.
func UpdatePost(db *sql.DB, post models.Post) error {
	_, err := db.Exec(`
		UPDATE posts
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, post.Title, post.Content, post.ID)

	return err
}

// DeletePost removes a post from the database.
func DeletePost(db *sql.DB, postID int) error {
	_, err := db.Exec(`
		DELETE FROM posts
		WHERE id = ?
	`, postID)

	return err
}

// queryPosts runs a query that returns multiple posts.
func queryPosts(db *sql.DB, query string, args ...any) ([]models.Post, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// scanPosts converts query rows into Post models.
func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
