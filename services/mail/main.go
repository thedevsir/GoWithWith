package mail

import (
	"os"
	"time"

	"github.com/Gommunity/GoWithWith/app/model"
	mailconfig "github.com/Gommunity/GoWithWith/config/mail"
	mailview "github.com/Gommunity/GoWithWith/resource/views/mail"

	jwt "github.com/dgrijalva/jwt-go"
	gomail "gopkg.in/gomail.v2"
)

func MakeEmailToken(action, username, email string, secret []byte) string {

	var token string
	var err error

	claims := model.EmailToken{
		action,
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err = tokenJWT.SignedString(secret); err != nil {
		panic(err)
	}

	return token
}

func SendVerficationMail(username, email, token string) {

	verifyBody := mailview.Verify{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := mailview.GenerateTemplate(verifyBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", verifyBody.EmailAddress)
	m.SetHeader("Subject", "Confirm your account")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := mailconfig.Driver.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendResetMail(username, email, token string) {

	forgotBody := mailview.Forgot{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := mailview.GenerateTemplate(forgotBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", forgotBody.EmailAddress)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := mailconfig.Driver.DialAndSend(m); err != nil {
		panic(err)
	}
}
