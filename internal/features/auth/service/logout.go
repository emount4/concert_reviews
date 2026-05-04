package auth_service

import "context"

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if s.authRepository == nil {
		return ErrAuthRepositoryNotConfigured
	}

	return s.authRepository.DeleteSession(ctx, refreshToken)
}
