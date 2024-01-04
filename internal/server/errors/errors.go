package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrNeededUpdate    = errors.New("needed update")
	ErrInvalidCurrency = errors.New("invalid currency")
)
