package service

import (
	"fmt"
	"jazz-network-api/model"
	"jazz-network-api/validation"
	"strings"
	"unicode"
)

var (
	organizerTypes = map[string]struct{}{
		"venue":    {},
		"festival": {},
	}
)

func validatePassword(password string) error {
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

func validateOrganizerInput(input model.OrganizerInput) error {
	var validationErr *model.ValidationError

	// type in ["venue", "festival"], required
	input.Type = strings.TrimSpace(input.Type)
	if input.Type == "" {
		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   "type",
			Message: "required",
		})
	}

	if _, found := organizerTypes[input.Type]; !found {
		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   "type",
			Message: "should be either 'venue' or 'festival'",
		})
	}

	// name required
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   "name",
			Message: "required",
		})
	}

	// country in existing real country, required
	input.Country = strings.TrimSpace(input.Country)

	// city required
	input.City = strings.TrimSpace(input.City)
	if input.City == "" {
		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   "city",
			Message: "required",
		})
	}

	// address required
	input.Address = strings.TrimSpace(input.Address)
	if input.Address == "" {
		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   "address",
			Message: "required",
		})
	}

	// websites valid format?
	for _, w := range input.Websites {
		// TODO: implement website format validation
		fmt.Println(w)
	}

	// emails valid format
	for _, e := range input.Emails {
		if !validation.IsValidEmail(e) {
			validationErr.Fields = append(validationErr.Fields, model.FieldError{
				Field:   "emails",
				Message: "wrong format",
			})
		}
	}

	// phone numbers valid format (maybe depending on country?), no letters
	for _, p := range input.PhoneNumbers {
		if !validation.IsValidPhoneNumber(e) {
			validationErr.Fields = append(validationErr.Fields, model.FieldError{
				Field:   "phone_numbers",
				Message: "wrong format",
			})
		}
	}
	// -90º <= latitude <= 90º, float
	// -180º <= longitude <= 180º, float
	return validationErr
}
