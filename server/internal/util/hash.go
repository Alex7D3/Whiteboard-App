package util

import (
	"golang.org/x/crypto/bcrypt"
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 12)
	return string(bytes), err
}

func CheckPasswordHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}

func HashToken(str string) string {
	sum := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}
