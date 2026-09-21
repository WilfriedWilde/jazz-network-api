package validation

import (
	"jazz-network-api/countries"
	"log"
	"net/mail"
	"net/url"

	"github.com/nyaruka/phonenumbers/v2"
)

func IsValidCountry(country string) bool {
	allCountries, err := countries.GetCountries()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if _, found := allCountries[country]; !found {
		return false
	}

	return true
}

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

func IsValidWebsite(w string) bool {
	u, err := url.ParseRequestURI(w)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
