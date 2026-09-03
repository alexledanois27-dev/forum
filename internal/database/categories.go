package database

import (
	"database/sql"

	"forum/internal/models"
)

// CreateCategory inserts a new category into the database.
func CreateCategory(db *sql.DB, category models.Category) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO categories (name, slug)
		VALUES (?, ?)
	`, category.Name, category.Slug)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllCategories returns all categories ordered by name.
func GetAllCategories(db *sql.DB) ([]models.Category, error) {
	rows, err := db.Query(`
		SELECT id, name, slug, created_at
		FROM categories
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryBySlug returns a category by its slug.
func GetCategoryBySlug(db *sql.DB, slug string) (*models.Category, error) {
	var category models.Category

	err := db.QueryRow(`
		SELECT id, name, slug, created_at
		FROM categories
		WHERE slug = ?
	`, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// GetCategoriesByPost returns all categories linked to a post.
func GetCategoriesByPost(db *sql.DB, postID int) ([]models.Category, error) {
	rows, err := db.Query(`
		SELECT c.id, c.name, c.slug, c.created_at
		FROM categories c
		INNER JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?
		ORDER BY c.name ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.CreatedAt,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// AddCategoryToPost links a category to a post.
func AddCategoryToPost(db *sql.DB, postID, categoryID int) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO post_categories (post_id, category_id)
		VALUES (?, ?)
	`, postID, categoryID)

	return err
}

// RemoveCategoryFromPost removes a category from a post.
func RemoveCategoryFromPost(db *sql.DB, postID, categoryID int) error {
	_, err := db.Exec(`
		DELETE FROM post_categories
		WHERE post_id = ? AND category_id = ?
	`, postID, categoryID)

	return err
}
