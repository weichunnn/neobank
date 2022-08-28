package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// return bcrypt hash string
// bcrypt + scrypt -> slow hashes
// sha1 + sha256 + sha512 -> fast hashes
// hash string format -> algo + cost + salt + hash
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// check password is correct when compared
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
