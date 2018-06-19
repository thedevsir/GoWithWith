package gomail

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/matcornic/hermes"
	gomail "gopkg.in/gomail.v2"
)

// Driver ...
var Driver *gomail.Dialer

// Initial ...
func Initial() {

	driverPort, _ := strconv.Atoi(os.Getenv("SMTPPort"))
	Driver = gomail.NewDialer(os.Getenv("SMTPHost"), driverPort, os.Getenv("SMTPUsername"), os.Getenv("SMTPPassword"))
}

// GenerateTemplate ...
func GenerateTemplate(email hermes.Email) (string, string) {

	year := strconv.Itoa(time.Now().Year())

	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: os.Getenv("EmailThemeName"),
			Link: os.Getenv("EmailThemeLink"),
			// Optional product logo
			Logo:        os.Getenv("EmailThemeLogo"),
			Copyright:   fmt.Sprintf(os.Getenv("EmailThemeCopyright"), year),
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	return emailBody, emailText
}
