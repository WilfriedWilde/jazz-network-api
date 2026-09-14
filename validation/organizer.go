package validation

import (
	"net/mail"
	"net/url"

	"github.com/nyaruka/phonenumbers/v2"
)

func IsValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func IsValidPhoneNumber(p string, country string) bool {
	num, err := phonenumbers.Parse(p, "")
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(num)
}

// TODO: implement website validation
func IsValidWebsite(w string) bool {
	u, err := url.ParseRequestURI(w)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
