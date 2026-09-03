package models

import "time"

type PostReaction struct {
	UserID    int
	PostID    int
	Value     int
	CreatedAt time.Time
}

type CommentReaction struct {
	UserID    int
	CommentID int
	Value     int
	CreatedAt time.Time
}
