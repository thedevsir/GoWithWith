package mail

import (
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

var Connection *gomail.Dialer

func Composer() {

	driverPort, err := strconv.Atoi(os.Getenv("SMTPPort"))

	if err != nil {
		panic(err)
	}

	Connection = gomail.NewDialer(
		os.Getenv("SMTPHost"),
		driverPort,
		os.Getenv("SMTPUsername"),
		os.Getenv("SMTPPassword"),
	)
}
