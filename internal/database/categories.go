package database

import (
	"database/sql"
	"fmt"

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

// initCategories seeds the categories expected by the existing templates.
// Existing names, slugs and IDs are preserved.
func initCategories(tx *sql.Tx) error {
	defaults := [][2]string{
		{"Annonces", "annonces"},
		{"Bugs et problèmes", "bugs-et-problemes"},
		{"Cartes & cartographie", "cartes-cartographie"},
		{"Changelog", "changelog"},
		{"Comptes rendus de parties", "comptes-rendus-de-parties"},
		{"Création de personnage", "creation-de-personnage"},
		{"Créations", "creations"},
		{"Discussion générale", "discussion-generale"},
		{"Discussions du lore", "discussions-du-lore"},
		{"Exploration", "exploration"},
		{"Fan arts", "fan-arts"},
		{"Fiches de personnages", "fiches-de-personnages"},
		{"Forum / Communauté", "forum-communaute"},
		{"Géographie", "geographie"},
		{"Histoire", "histoire"},
		{"Hors sujet", "hors-sujet"},
		{"Illustrations", "illustrations"},
		{"JDR", "jdr"},
		{"Lore", "lore"},
		{"Magie & cosmologie", "magie-cosmologie"},
		{"Mises à jour du lore", "mises-a-jour-du-lore"},
		{"Nature", "nature"},
		{"Parties & campagnes", "parties-campagnes"},
		{"Personnages", "personnages"},
		{"Peuples & factions", "peuples-factions"},
		{"Présentations", "presentations"},
		{"Races", "races"},
		{"Religions & mythes", "religions-mythes"},
		{"Règlement", "reglement"},
		{"Règles du JDR", "regles-du-jdr"},
		{"Suggestions pour le forum", "suggestions-pour-le-forum"},
		{"Technologies & artefacts", "technologies-artefacts"},
		{"Textes et nouvelles", "textes-et-nouvelles"},
	}
	for _, category := range defaults {
		// The forum parent has two legacy aliases accepted by the templates.
		if category[0] == "Forum / Communauté" {
			var count int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM categories WHERE name IN ('Forum', 'Communauté')`).Scan(&count); err != nil {
				return err
			}
			if count > 0 {
				continue
			}
		}
		if _, err := tx.Exec(`INSERT INTO categories (name, slug) VALUES (?, ?) ON CONFLICT(name) DO NOTHING`, category[0], category[1]); err != nil {
			return fmt.Errorf("category %q: %w", category[0], err)
		}
	}
	return nil
}
