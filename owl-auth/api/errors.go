package api

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotFound        = errors.New("user not found")
	ErrTooManyAttempts = errors.New("too many attempts. please reset password")
)
