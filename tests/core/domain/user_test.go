package domain_test

import (
	"errors"
	"testing"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

func TestUserValidateSuccess(t *testing.T) {
	user := domain.User{
		Email:        "user@example.com",
		PasswordHash: "hash",
		Username:     "user_name",
		RoleID:       1,
	}
	if err := user.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestUserValidateRejectsUppercaseEmail(t *testing.T) {
	user := domain.User{
		Email:        "User@Example.com",
		PasswordHash: "hash",
		Username:     "user_name",
		RoleID:       1,
	}
	err := user.Validate()
	if err == nil {
		t.Fatal("Validate() expected error for uppercase email, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatal("Validate() expected invalid argument error")
	}
}

func TestUserValidateRejectsInvalidUsername(t *testing.T) {
	user := domain.User{
		Email:        "user@example.com",
		PasswordHash: "hash",
		Username:     "bad name",
		RoleID:       1,
	}
	if err := user.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid username, got nil")
	}
}

func TestUserApplyPatchTrimsUsername(t *testing.T) {
	user := domain.User{
		Email:        "user@example.com",
		PasswordHash: "hash",
		Username:     "user_name",
		RoleID:       1,
	}
	name := "  new_user  "
	patch := domain.UserPatch{Username: core_types.Nullable[string]{Value: &name, Set: true}}

	if err := user.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}
	if user.Username != "new_user" {
		t.Fatalf("ApplyPatch() expected trimmed username, got %q", user.Username)
	}
}

func TestUserApplyPatchRejectsNullUsername(t *testing.T) {
	user := domain.User{
		Email:        "user@example.com",
		PasswordHash: "hash",
		Username:     "user_name",
		RoleID:       1,
	}
	patch := domain.UserPatch{Username: core_types.Nullable[string]{Value: nil, Set: true}}
	err := user.ApplyPatch(patch)
	if err == nil {
		t.Fatal("ApplyPatch() expected error for null username, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatal("ApplyPatch() expected invalid argument error")
	}
}

func TestUserPatchIsEmpty(t *testing.T) {
	patch := domain.UserPatch{}
	if !patch.IsEmpty() {
		t.Fatal("IsEmpty() expected true for empty patch")
	}
	name := "user"
	patch.Username = core_types.Nullable[string]{Value: &name, Set: true}
	if patch.IsEmpty() {
		t.Fatal("IsEmpty() expected false when patch has values")
	}
}
