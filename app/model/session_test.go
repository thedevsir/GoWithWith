package model

import (
	"testing"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/Gommunity/GoWithWith/services/test"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

var sessionCollection *mongodm.Model

func sessionBeforeTest() {

	utility.LoadEnvironmentVariables("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot(map[string]mongodm.IDocumentBase{
		"sessions": &Session{},
	})

	sessionCollection = database.Connection.Model(SessionCollection)
	sessionCollection.RemoveAll(nil)
}

func sessionAfterTest() {
	sessionCollection.RemoveAll(nil)
}

func TestSession(t *testing.T) {

	sessionBeforeTest()

	var SID string
	var sess string

	session := &Session{
		IP:        "127.0.0.1",
		UserID:    "1111111111",
		UserAgent: ":::USER-AGENT:::",
	}

	t.Run("CreateSession", func(t *testing.T) {
		assert.NotPanics(t, func() {
			SID, sess = session.Create()
		})
	})

	t.Run("SessionFindByID", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			session, err := SessionFindByID(SID)
			assert.NoError(t, err)
			assert.IsType(t, Session{}, session)
		})

		t.Run("SessionNotFound", func(t *testing.T) {
			_, err := SessionFindByID("54759eb3c090d83494e2d804")
			assert.Error(t, err)
			assert.Equal(t, "Session was not found", err.Error())
		})
	})

	t.Run("SessionFindByCredentials", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			err := SessionFindByCredentials(sess, SID)
			assert.NoError(t, err)
		})

		t.Run("FakeSession", func(t *testing.T) {
			err := SessionFindByCredentials("fake", SID)
			assert.Error(t, err)
			assert.Equal(t, "Credentials are invalid", err.Error())
		})
	})

	t.Run("SessionUpdateLastActivity", func(t *testing.T) {
		assert.NotPanics(t, func() {
			SessionUpdateLastActivity(SID)
		})
	})

	t.Run("GetUserSessions", func(t *testing.T) {
		var sessions paginate.Paginate
		sessions, err := GetUserSessions("fake", 1, 10)
		assert.NoError(t, err)
		assert.IsType(t, paginate.Paginate{}, sessions)
	})

	t.Run("DeleteSession", func(t *testing.T) {
		assert.NotPanics(t, func() { DeleteSession(SID) })
	})

	sessionAfterTest()
}
