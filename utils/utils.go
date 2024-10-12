package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashBCrypt(data []byte) ([]byte, error) {
	hashedData, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)

	return hashedData, err
}
