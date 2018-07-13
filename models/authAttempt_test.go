package models

import (
	"testing"

	"github.com/Gommunity/GoWithWith/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

var dbAuthAttempt *mongodm.Model

func authAttemptBeforeTest() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic(":main:init: ErrorLoading.EnvFile")
	}

	Models := make(map[string]mongodm.IDocumentBase)
	Models["authAttempts"] = &AuthAttempt{}

	database.Initial(Models, true)

	dbAuthAttempt = database.Connection.Model("AuthAttempt")
	dbAuthAttempt.RemoveAll(nil)
}

func authAttemptAfterTest() {
	dbAuthAttempt.RemoveAll(nil)
}

func TestAuthAttempt(t *testing.T) {

	authAttemptBeforeTest()

	ad := AbuseDetected{
		MaxIP:            "2",
		MaxIPAndUsername: "5",
	}

	t.Run("AttemptCreateSuccess", func(t *testing.T) {
		err := AttemptCreate("::1", "Irani")
		assert.Nil(t, err)
	})

	t.Run("AbuseDetectedSuccess", func(t *testing.T) {
		err := ad.Check("::1", "Irani")
		assert.Nil(t, err)
	})

	t.Run("AbuseDetectedFailure", func(t *testing.T) {
		AttemptCreate("::1", "Irani")
		err := ad.Check("::1", "Irani")
		assert.Error(t, err)
	})

	authAttemptAfterTest()
}
