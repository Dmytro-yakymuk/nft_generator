package mail

// ensures that MockMailer implements Mailer.
var _ Mailer = (*MockMailer)(nil)

// MockMailer is mock implementation of Mailer.
type MockMailer struct{}

// FromAddress returns address of sender from.
func (sender *MockMailer) FromAddress() Address {
	return Address{
		Name:    "",
		Address: "",
	}
}

// Send sends email message to the given recipient.
func (sender *MockMailer) Send(msg *Email) error {
	return nil
}
