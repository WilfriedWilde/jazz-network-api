package validation

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort = errors.New("service: password too short")
	ErrPasswordTooWeak  = errors.New("service: password too weak")
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasLetter, hasNumber bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsNumber(r):
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
}
