package viewmail

import (
	"fmt"
	"os"

	"github.com/matcornic/hermes"
)

type Forgot struct {
	Username     string
	EmailAddress string
	Token        string
}

func (r *Forgot) Name() string {
	return "forgot"
}

func (r *Forgot) Email() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: r.Username,
			Intros: []string{
				"You have received this email because a password reset request for \"" + r.EmailAddress + "\" account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  fmt.Sprintf(os.Getenv("EmailResetLink"), r.Token),
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
}
