package mail

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/matcornic/hermes"
)

func GenerateTemplate(email hermes.Email) (string, string) {

	year := strconv.Itoa(time.Now().Year())
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        os.Getenv("EmailThemeName"),
			Link:        os.Getenv("EmailThemeLink"),
			Logo:        os.Getenv("EmailThemeLogo"),
			Copyright:   fmt.Sprintf(os.Getenv("EmailThemeCopyright"), year),
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err)
	}

	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		panic(err)
	}

	return emailBody, emailText
}
