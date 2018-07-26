package mail

import (
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

var Driver *gomail.Dialer

func Initial() {
	driverPort, _ := strconv.Atoi(os.Getenv("SMTPPort"))
	Driver = gomail.NewDialer(os.Getenv("SMTPHost"), driverPort, os.Getenv("SMTPUsername"), os.Getenv("SMTPPassword"))
}
