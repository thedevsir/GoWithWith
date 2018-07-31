package mail

import (
	"os"
	"testing"

	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
	gomail "gopkg.in/gomail.v2"
)

func TestMailConnection(t *testing.T) {

	t.Run("WithOutConfigInfo", func(t *testing.T) {
		assert.Panics(t, func() { Composer() })
	})

	t.Run("SendTestEmail", func(t *testing.T) {

		utility.LoadEnvironmentVariables("../../.env")
		Composer()
		m := gomail.NewMessage()
		m.SetHeader("From", os.Getenv("EmailFrom"))
		m.SetHeader("To", "freshmanlimited@gmail.com")
		m.SetHeader("Subject", "SendTestEmail")
		m.SetBody("text/plain", "...")

		assert.NoError(t, Connection.DialAndSend(m))
	})
}
