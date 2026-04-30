package auth_service

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
	ValidatePassword(password, hash string) bool
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}

func (h *SHA1Hasher) ValidatePassword(password, hash string) bool {
	return hash == h.Hash(password)
}
