package validation

import (
	"net/mail"
	"net/url"
)

func IsValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

// TODO: implement phone number validation -> nyaruka/phonenumbers
func IsValidPhoneNumber(p string) bool {
	return true
}

// TODO: implement website validation
func IsValidWebsite(w string) bool {
	u, err := url.ParseRequestURI(w)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
