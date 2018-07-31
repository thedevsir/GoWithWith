package mail

import (
	"testing"

	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
)

func TestTemplates(t *testing.T) {

	utility.LoadEnvironmentVariables("../../../.env")

	t.Run("Verify", func(t *testing.T) {

		verifyBody := Verify{
			Username:     "fakeUser",
			EmailAddress: "fakeEmail",
			Token:        "fakeToken",
		}

		assert.Equal(t, "verify", verifyBody.Name())
		assert.NotPanics(t, func() { GenerateTemplate(verifyBody.Email()) })
	})

	t.Run("Forgot", func(t *testing.T) {

		forgotBody := Forgot{
			Username:     "fakeUser",
			EmailAddress: "fakeEmail",
			Token:        "fakeToken",
		}

		assert.Equal(t, "forgot", forgotBody.Name())
		assert.NotPanics(t, func() { GenerateTemplate(forgotBody.Email()) })
	})
}
