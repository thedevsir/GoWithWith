package mail

import (
	"testing"

	"github.com/Gommunity/GoWithWith/config/mail"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
)

func TestMakeTokenForEmails(t *testing.T) {

	assert.NotPanics(t, func() {
		MakeTokenForEmails("verify", "username", "@", []byte("secret"))
	})
}

func TestSendMail(t *testing.T) {

	utility.LoadEnvironmentVariables("../../.env")
	mail.Composer()

	t.Run("SendVerficationMail", func(t *testing.T) {
		assert.NotPanics(t, func() { SendVerficationMail("username", "freshmanlimited@gmail.com", "token") })
	})

	t.Run("SendResetMail", func(t *testing.T) {
		assert.NotPanics(t, func() { SendResetMail("username", "freshmanlimited@gmail.com", "token") })
	})
}
