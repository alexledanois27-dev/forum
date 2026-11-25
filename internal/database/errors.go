package database

import (
	"errors"
	"github.com/mattn/go-sqlite3"
)

// ErrInvalidReactionValue is returned when a reaction value is not 1 or -1.
var ErrInvalidReactionValue = errors.New("reaction value must be 1 or -1")

// IsUniqueViolation distinguishes unique constraint errors from technical failures.
func IsUniqueViolation(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
