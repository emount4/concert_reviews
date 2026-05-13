package core_errors

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrInvalidArgument         = errors.New("invalid argument")
	ErrConflict                = errors.New("conflict")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrRepositoryNotConfigured = errors.New("not configured")
	ErrTxManagerNotConfigured  = errors.New("transaction manager is not configured")
	ErrForbidden               = errors.New("forbidden")
)
