package utils

import "golang.org/x/crypto/bcrypt"

func ComparePassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
