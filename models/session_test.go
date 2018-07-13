package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Gommunity/GoWithWith/database"
	"github.com/joho/godotenv"
	"github.com/zebresel-com/mongodm"
)

var dbSession *mongodm.Model

func sessionBeforeTest() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic(":main:init: ErrorLoading.EnvFile")
	}

	Models := make(map[string]mongodm.IDocumentBase)
	Models["session"] = &Session{}

	database.Initial(Models, true)

	dbSession = database.Connection.Model("session")
	dbSession.RemoveAll(nil)
}

func sessionAfterTest() {
	dbSession.RemoveAll(nil)
}

func TestSession(t *testing.T) {

	sessionBeforeTest()

	id, session := CreateSession("Irani", "::1", "userAgent")

	t.Run("CreateSession", func(t *testing.T) {
		assert.IsType(t, String, id)
		assert.IsType(t, String, session)
	})

	t.Run("SessionFindByID", func(t *testing.T) {
		session, err := SessionFindByID(id)
		assert.Nil(t, err)
		assert.IsType(t, Session{}, session)
	})

	t.Run("FindByCredentials", func(t *testing.T) {
		err := SessionFindByCredentials(session, id)
		err1 := SessionFindByCredentials("fake", id)
		assert.Nil(t, err)
		assert.Error(t, err1)
	})

	t.Run("UpdateLastActive", func(t *testing.T) {
		err := UpdateLastActive(id)
		assert.Nil(t, err)
	})

	t.Run("GetUserSessions", func(t *testing.T) {
		sessions, err := GetUserSessions(id, 1, 10)
		assert.Nil(t, err)
		assert.IsType(t, Pagination{}, sessions)
	})

	DeleteSession(id)

	sessionAfterTest()
}

func TestCreateJWToken(t *testing.T) {

	jwt := CreateJWToken("session", "SID", "Irani", "userID", []byte("Secret"))

	assert.IsType(t, String, jwt)
}
