package mail

import (
	"os"
	"time"

	Mail "github.com/Gommunity/GoWithWith/config/mail"
	"github.com/Gommunity/GoWithWith/resource/views/mail"
	jwt "github.com/dgrijalva/jwt-go"
	gomail "gopkg.in/gomail.v2"
)

type EmailToken struct {
	Action   string
	Email    string
	Username string
	jwt.StandardClaims
}

func MakeTokenForEmails(action, username, email string, secret []byte) string {

	claims := EmailToken{
		action,
		email,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJWT.SignedString(secret)

	if err != nil {
		panic(err)
	}

	return tokenString
}

func SendVerficationMail(username, email, token string) {

	verifyBody := &mail.Verify{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := mail.GenerateTemplate(verifyBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", verifyBody.EmailAddress)
	m.SetHeader("Subject", "Confirm your account")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := Mail.Connection.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendResetMail(username, email, token string) {

	forgotBody := mail.Forgot{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := mail.GenerateTemplate(forgotBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", forgotBody.EmailAddress)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := Mail.Connection.DialAndSend(m); err != nil {
		panic(err)
	}
}
