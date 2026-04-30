// Задание: синхронизировать структуру `User` с миграциями (UUID, nullable поля, теги)
package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type User struct {
	ID uuid.UUID `json:"id" db:"user_id"`

	Email            string  `json:"email" db:"email"`
	PasswordHash     string  `json:"password_hash" db:"password_hash"`
	TelegramID       *int64  `json:"telegram_id,omitempty" db:"tg_id"`
	TelegramUsername *string `json:"telegram_username,omitempty" db:"tg_username"`
	RoleID           int     `json:"role_id" db:"role_id"`
	Username         string  `json:"username" db:"username"`
	Bio              *string `json:"bio,omitempty" db:"bio"`

	AvatarURL *string `json:"avatar_url,omitempty" db:"avatar_url"`
	BannerURL *string `json:"banner_url,omitempty" db:"banner_url"`

	IsEmailVerified bool       `json:"is_email_verified" db:"is_email_verified"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsBanned        bool       `json:"is_banned" db:"is_banned"`
	BannedByUserID  *uuid.UUID `json:"banned_by_user_id,omitempty" db:"banned_by_user_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(username, email string) User {
	return User{
		Username:        username,
		Email:           email,
		RoleID:          1,
		IsEmailVerified: false,
		IsActive:        true,
		IsBanned:        false,
	}
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("%w: email is required", core_errors.ErrInvalidArgument)
	}
	if utf8.RuneCountInString(u.Email) > 255 {
		return fmt.Errorf("%w: email max length is 255", core_errors.ErrInvalidArgument)
	}
	if u.Email != strings.ToLower(u.Email) {
		return fmt.Errorf("%w: email must be lowercase", core_errors.ErrInvalidArgument)
	}
	if strings.ContainsAny(u.Email, " \t\n\r") {
		return fmt.Errorf("%w: email must not contain whitespace", core_errors.ErrInvalidArgument)
	}
	if at := strings.Index(u.Email, "@"); at <= 0 {
		return fmt.Errorf("%w: email must contain '@' after first character", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return fmt.Errorf("%w: password_hash is required", core_errors.ErrInvalidArgument)
	}
	if utf8.RuneCountInString(u.PasswordHash) > 255 {
		return fmt.Errorf("%w: password_hash max length is 255", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(u.Username) == "" {
		return fmt.Errorf("%w: username is required", core_errors.ErrInvalidArgument)
	}
	usernameLen := utf8.RuneCountInString(u.Username)
	if usernameLen < 3 {
		return fmt.Errorf("%w: username min length is 3", core_errors.ErrInvalidArgument)
	}
	if usernameLen > 50 {
		return fmt.Errorf("%w: username max length is 50", core_errors.ErrInvalidArgument)
	}
	if !usernameRegexp.MatchString(u.Username) {
		return fmt.Errorf("%w: username must match ^[a-zA-Z0-9_]+$", core_errors.ErrInvalidArgument)
	}

	if u.RoleID <= 0 {
		return fmt.Errorf("%w: role_id must be positive", core_errors.ErrInvalidArgument)
	}

	if u.TelegramUsername != nil && utf8.RuneCountInString(*u.TelegramUsername) > 100 {
		return fmt.Errorf("%w: tg_username max length is 100", core_errors.ErrInvalidArgument)
	}

	if u.Bio != nil && utf8.RuneCountInString(*u.Bio) > 500 {
		return fmt.Errorf("%w: bio max length is 500", core_errors.ErrInvalidArgument)
	}

	if u.AvatarURL != nil {
		if strings.TrimSpace(*u.AvatarURL) == "" {
			return fmt.Errorf("%w: avatar_url must not be blank", core_errors.ErrInvalidArgument)
		}
		if utf8.RuneCountInString(*u.AvatarURL) > 2048 {
			return fmt.Errorf("%w: avatar_url max length is 2048", core_errors.ErrInvalidArgument)
		}
	}

	if u.BannerURL != nil {
		if strings.TrimSpace(*u.BannerURL) == "" {
			return fmt.Errorf("%w: banner_url must not be blank", core_errors.ErrInvalidArgument)
		}
		if utf8.RuneCountInString(*u.BannerURL) > 2048 {
			return fmt.Errorf("%w: banner_url max length is 2048", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
