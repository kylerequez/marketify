package utils

import "golang.org/x/crypto/bcrypt"

func ComparePassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func EncryptPassword(password string, cost int) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}

	return hashed, nil
}
