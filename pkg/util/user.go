package util

import (
	"crypto/sha512"
	"encoding/base64"
	"user_crud/internal/pkg/setting"
)

func HashPassword(password string) string {

	salt := []byte(setting.SECRET_KEY)

	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a base64 encoded string
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

func IsPasswordMatch(hashed, password string) bool {
	return hashed == HashPassword(password)
}
