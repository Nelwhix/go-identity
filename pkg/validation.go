package pkg

import (
	"fmt"
	"regexp"
)

func StrictPasswordValidation(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	uppercase := regexp.MustCompile(`[A-Z]`)
	if !uppercase.MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	specialChar := regexp.MustCompile(`[!@#~$%^&*()_+\-={}|[\]\\:";'<>?,./]`)
	if !specialChar.MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	number := regexp.MustCompile(`[0-9]`)
	if !number.MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	return nil
}
