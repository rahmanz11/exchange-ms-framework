package formaterror

import (
	"errors"
	"strings"
)

// FormatError formats the error message
func FormatError(err string) error {

	if strings.Contains(err, "account number") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "credential") {
		return errors.New("Incorrect Password")
	}

	return errors.New("Incorrect Details")
}
