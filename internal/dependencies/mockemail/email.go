package mockemail

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) Send(email *domain.Email) *domain.EmailError {
	fmt.Printf("sending email from %s to %s\n", email.FromAddress, email.ToAddress)
	fmt.Printf("subject: %s\n", email.Subject)
	fmt.Printf("contents: %s\n", email.Content)
	return nil
}
