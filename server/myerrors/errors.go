package myerrors

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidInput    = errors.New("invalid input")
	ErrCartEmpty       = errors.New("cart is empty")
	ErrEmailTaken      = errors.New("email already taken")
	ErrSessionExpired  = errors.New("session expired")
	ErrInternal        = errors.New("internal error")
)
