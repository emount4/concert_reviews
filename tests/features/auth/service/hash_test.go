package auth_service_test

import (
	"testing"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func TestSHA1HasherValidatePassword(t *testing.T) {
	hasher := auth_service.NewSHA1Hasher("salt")
	hash := hasher.Hash("pass")

	if !hasher.ValidatePassword("pass", hash) {
		t.Fatal("expected password to validate")
	}
	if hasher.ValidatePassword("wrong", hash) {
		t.Fatal("expected password validation to fail")
	}
}
