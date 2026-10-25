package mail

import (
	"strings"
)

// Normalize brings the email to UpperCase.
func Normalize(email string) string {
	return strings.ToUpper(email)
}
