package models

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
