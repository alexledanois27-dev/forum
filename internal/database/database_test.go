package database

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"path/filepath"
	"testing"
	"time"
)

func TestDatabaseAudit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "forum.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { db.Close() }()
	cats, err := GetAllCategories(db)
	if err != nil || len(cats) != 33 {
		t.Fatalf("categories: %d, %v", len(cats), err)
	}
	a, err := CreateUser(db, models.User{Username: "alice", Email: "alice@example.test", PasswordHash: "hash"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := CreateUser(db, models.User{Username: "bob", Email: "bob@example.test", PasswordHash: "hash"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateUser(db, models.User{Username: "ALICE", Email: "other@example.test", PasswordHash: "hash"}); !IsUniqueViolation(err) {
		t.Fatalf("unique error: %v", err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	for _, id := range []string{"old", "new"} {
		if err := CreateSession(db, models.Session{ID: id, UserID: a, ExpiresAt: now.Add(time.Hour)}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := GetSession(db, "old"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("old session: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO sessions(id,user_id,expires_at) VALUES ('duplicate',?,?)`, a, now.Add(time.Hour)); !IsUniqueViolation(err) {
		t.Fatalf("session uniqueness: %v", err)
	}
	if _, err := GetValidSession(db, "new", now); err != nil {
		t.Fatal(err)
	}
	if _, err := GetValidSession(db, "new", now.Add(time.Hour)); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expired session: %v", err)
	}
	post := models.Post{UserID: a, Title: "Test", Content: "Content"}
	p, err := CreatePost(db, post, []int{int(cats[0].ID), int(cats[1].ID), int(cats[0].ID)})
	if err != nil {
		t.Fatal(err)
	}
	linked, err := GetCategoriesByPost(db, int(p))
	if err != nil || len(linked) != 2 {
		t.Fatalf("links: %v, %v", linked, err)
	}
	if _, err := CreatePost(db, post, []int{-1}); err == nil {
		t.Fatal("invalid foreign key accepted")
	}
	all, err := GetAllPosts(db)
	if err != nil || len(all) != 1 {
		t.Fatalf("rollback: %v, %v", all, err)
	}
	if err := SetPostReaction(db, int(a), int(p), 1); err != nil {
		t.Fatal(err)
	}
	if err := SetPostReaction(db, int(b), int(p), -1); err != nil {
		t.Fatal(err)
	}
	liked, err := GetLikedPostsByUser(db, int(a))
	if err != nil || len(liked) != 1 || liked[0].ID != p {
		t.Fatalf("liked: %v, %v", liked, err)
	}
	liked, err = GetLikedPostsByUser(db, int(b))
	if err != nil || len(liked) != 0 {
		t.Fatalf("isolation: %v, %v", liked, err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	again, err := GetAllCategories(db)
	if err != nil || len(again) != len(cats) || again[0].ID != cats[0].ID {
		t.Fatalf("seed not idempotent: %v", err)
	}
}

func TestLegacySessionMigration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	id, err := CreateUser(db, models.User{Username: "legacy", Email: "legacy@example.test", PasswordHash: "hash"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`DROP INDEX idx_sessions_unique_user_id`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES ('old',?,'2099-01-01 00:00:00+00:00','2020-01-01'),('latest',?,'2099-01-01 00:00:00+00:00','2021-01-01')`, id, id); err != nil {
		t.Fatal(err)
	}
	db.Close()
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := GetSession(db, "old"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("old session: %v", err)
	}
	if _, err := GetValidSession(db, "latest", time.Now()); err != nil {
		t.Fatal(err)
	}
}
