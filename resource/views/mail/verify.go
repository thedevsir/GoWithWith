package mail

import (
	"fmt"
	"os"

	"github.com/matcornic/hermes"
)

type Verify struct {
	Username     string
	EmailAddress string
	Token        string
}

func (w *Verify) Name() string {
	return "verify"
}

func (w *Verify) Email() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: w.Username,
			Intros: []string{
				"Welcome to " + os.Getenv("EmailAppName") + "! We're very excited to have you on board.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Username", Value: w.Username},
				{Key: "Email", Value: w.EmailAddress},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with " + os.Getenv("EmailAppName") + ", please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: fmt.Sprintf(os.Getenv("EmailVerifyLink"), w.Token),
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
