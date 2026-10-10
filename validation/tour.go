package validation

import (
	"time"
)

// IsValidDate checks if date has format YYYY-MM-DD, is a valid calendar date and is situated in the future
func IsValidDate(date string) (bool, string) {
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return false, err.Error()
	}

	if d.Before(time.Now()) {
		return false, "should be in the future"
	}

	return true, ""
}
