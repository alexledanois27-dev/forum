package utils

import "golang.org/x/crypto/bcrypt"

const bcryptCost = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
