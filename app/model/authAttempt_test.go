package model

import (
	"testing"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/test"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

var authAttemptCollection *mongodm.Model

func authAttemptBeforeTest() {

	utility.LoadEnvironmentVariables("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot(map[string]mongodm.IDocumentBase{
		"authAttempts": &AuthAttempt{},
	})

	authAttemptCollection = database.Connection.Model(AuthAttemptCollection)
	authAttemptCollection.RemoveAll(nil)
}

func authAttemptAfterTest() {
	authAttemptCollection.RemoveAll(nil)
}

func TestAuthAttempt(t *testing.T) {

	authAttemptBeforeTest()

	attempt := &AuthAttempt{
		IP:       "127.0.0.1",
		Username: "Irani",
	}

	t.Run("CreateAttempt", func(t *testing.T) {
		for i := 0; i <= 5; i++ {
			assert.NotPanics(t, func() { attempt.Create() })
		}
	})

	t.Run("MaximumAttemptsNotReached", func(t *testing.T) {

		abuse := &AbuseDetected{
			MaxIP:            "10",
			MaxIPAndUsername: "10",
		}

		err := abuse.Check(attempt.IP, attempt.Username)
		assert.Nil(t, err)
	})

	t.Run("MaximumAttemptsReached", func(t *testing.T) {

		abuse := &AbuseDetected{
			MaxIP:            "3",
			MaxIPAndUsername: "3",
		}

		err := abuse.Check(attempt.IP, attempt.Username)
		assert.Error(t, err)
		assert.Equal(t, "Maximum number of auth attempts reached", err.Error())
	})

	authAttemptAfterTest()
}
