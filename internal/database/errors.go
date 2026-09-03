package database

import "errors"

// ErrInvalidReactionValue is returned when a reaction value is not 1 or -1.
var ErrInvalidReactionValue = errors.New("reaction value must be 1 or -1")
