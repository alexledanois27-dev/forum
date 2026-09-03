package utils

import (
	"errors"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 254 {
		return errors.New("email is too long")
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return errors.New("invalid email address")
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if username != strings.TrimSpace(username) {
		return errors.New("username cannot contain surrounding spaces")
	}

	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must contain between 3 and 30 characters")
	}

	for _, character := range username {
		valid := character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			character == '_' ||
			character == '-'

		if !valid {
			return errors.New("username contains an invalid character")
		}
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must contain at least 8 characters")
	}

	if len(password) > 72 {
		return errors.New("password cannot exceed 72 bytes")
	}

	return nil
}
