package validation

import (
	"jazz-network-api/countries"
	"jazz-network-api/model"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/nyaruka/phonenumbers/v2"
)

var (
	organizerTypes = map[string]struct{}{
		"venue":    {},
		"festival": {},
	}
)

func ValidateOrganizerInput(input model.OrganizerInput) error {
	var validationErr *model.ValidationError

	addFieldErr := func(field, message string) {
		if validationErr == nil {
			validationErr = &model.ValidationError{
				Status: http.StatusBadRequest,
				Code:   "validation_error",
			}
		}

		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   field,
			Message: message,
		})
	}

	// type in ["venue", "festival"], required
	input.Type = strings.TrimSpace(input.Type)
	if input.Type == "" {
		addFieldErr("type", "required")
	} else if _, found := organizerTypes[input.Type]; !found {
		addFieldErr("type", "should be either 'venue' or 'festival'")
	}

	// name required
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		addFieldErr("name", "required")
	}

	// country in existing real country, required
	input.Country = strings.TrimSpace(input.Country)
	if input.Country == "" {
		addFieldErr("country", "required")
	} else if !isValidCountry(input.Country) {
		addFieldErr("country", "invalid country")
	}

	// city required
	input.City = strings.TrimSpace(input.City)
	if input.City == "" {
		addFieldErr("city", "required")
	}

	// address required
	input.Address = strings.TrimSpace(input.Address)
	if input.Address == "" {
		addFieldErr("address", "required")
	}

	// websites valid format
	for _, w := range input.Websites {
		if !isValidWebsite(w) {
			addFieldErr("websites", "wrong format")
			break
		}
	}

	// emails valid format
	for _, e := range input.Emails {
		if !IsValidEmail(e) {
			addFieldErr("emails", "wrong format")
			break
		}
	}

	// phone numbers valid format
	for _, p := range input.PhoneNumbers {
		if !isValidPhoneNumber(p, input.Country) {
			addFieldErr("phone_numbers", "wrong format")
			break
		}
	}

	// -90º <= latitude <= 90º, required
	if input.Latitude == nil {
		addFieldErr("latitude", "required")
	} else if *input.Latitude < -90. || *input.Latitude > 90. {
		addFieldErr("latitude", "should be in range [-90, 90]")
	}

	// -180º <= longitude <= 180º, required
	if input.Longitude == nil {
		addFieldErr("longitude", "required")
	} else if *input.Longitude < -180. || *input.Longitude > 180. {
		addFieldErr("longitude", "should be in range [-180, 180]")
	}

	// public visibility required
	if input.PublicVisibility == nil {
		addFieldErr("public_visibility", "required")
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}

func ValidateOrganizerUpdate(update model.OrganizerUpdate) error {
	var validationErr *model.ValidationError

	addFieldErr := func(field, message string) {
		if validationErr == nil {
			validationErr = &model.ValidationError{
				Status: http.StatusBadRequest,
				Code:   "validation_error",
			}
		}

		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   field,
			Message: message,
		})
	}

	// type in ["venue", "festival"]
	if update.Type != nil {
		*update.Type = strings.TrimSpace(*update.Type)
		if _, found := organizerTypes[*update.Type]; !found {
			addFieldErr("type", "should be either 'venue' or 'festival'")
		}
	}

	// name nonempty
	if update.Name != nil {
		*update.Name = strings.TrimSpace(*update.Name)
		if *update.Name == "" {
			addFieldErr("name", "nonempty")
		}
	}

	// country in existing real country
	if update.Country != nil {
		if !isValidCountry(*update.Country) {
			addFieldErr("country", "invalid country")
		}
	}

	// city nonempty
	if update.City != nil {
		*update.City = strings.TrimSpace(*update.City)
		if *update.City == "" {
			addFieldErr("city", "nonempty")
		}
	}

	// address nonempty
	if update.Address != nil {
		*update.Address = strings.TrimSpace(*update.Address)
		if *update.Address == "" {
			addFieldErr("address", "nonempty")
		}
	}

	// websites valid format
	for _, w := range update.Websites {
		if !isValidWebsite(w) {
			addFieldErr("websites", "wrong format")
			break
		}
	}

	// emails valid format
	for _, e := range update.Emails {
		if !IsValidEmail(e) {
			addFieldErr("emails", "wrong format")
			break
		}
	}

	// phone numbers valid format
	if update.Country != nil {
		for _, p := range update.PhoneNumbers {
			if !isValidPhoneNumber(p, *update.Country) {
				addFieldErr("phone_numbers", "wrong format")
				break
			}
		}
	}

	// -90º <= latitude <= 90º
	if update.Latitude != nil {
		if *update.Latitude < -90. || *update.Latitude > 90. {
			addFieldErr("latitude", "should be in range [-90, 90]")
		}
	}

	// -180º <= longitude <= 180º
	if update.Longitude != nil {
		if *update.Longitude < -180. || *update.Longitude > 180. {
			addFieldErr("longitude", "should be in range [-180, 180]")
		}
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}

func isValidCountry(country string) bool {
	allCountries, err := countries.GetCountries()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	byCode := make(map[string]struct{}, len(allCountries))
	for _, c := range allCountries {
		byCode[c.Code] = struct{}{}
	}

	if _, found := byCode[country]; !found {
		return false
	}

	return true
}

func IsValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func isValidPhoneNumber(p string, country string) bool {
	num, err := phonenumbers.Parse(p, "")
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(num)
}

func isValidWebsite(w string) bool {
	u, err := url.ParseRequestURI(w)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
