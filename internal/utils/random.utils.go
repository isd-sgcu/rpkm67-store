package utils

import (
	"crypto/rand"
	"encoding/base64"
)

type Utils interface {
	GenerateRandomString(length int) (string, error)
}

type randomUtils struct{}

func NewRandomUtils() Utils {
	return &randomUtils{}
}

func (u *randomUtils) GenerateRandomString(length int) (string, error) {
	numBytes := (length * 6) / 8

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	return randomString[:length], nil
}
