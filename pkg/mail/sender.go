package mail

import (
	"net/mail"
)

// Address is alias of net/mail.Address.
type Address = mail.Address

// Mailer sends emails.
type Mailer interface {
	FromAddress() Address
	Send(msg *Email) error
}
