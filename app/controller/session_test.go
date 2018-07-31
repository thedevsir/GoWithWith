package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/stretchr/testify/assert"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/auth"
	j "github.com/Gommunity/GoWithWith/services/jwt"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/Gommunity/GoWithWith/services/test"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/zebresel-com/mongodm"
)

var sessionCollection *mongodm.Model

func sessionBeforeTest() {

	utility.LoadEnvironmentVariables("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot(map[string]mongodm.IDocumentBase{
		"sessions": &model.Session{},
	})

	sessionCollection = database.Connection.Model(model.SessionCollection)
	sessionCollection.RemoveAll(nil)
}

func sessionAfterTest() {
	sessionCollection.RemoveAll(nil)
}

func TestSessions(t *testing.T) {

	sessionBeforeTest()

	token := &auth.JWTUserToken{
		Session:    "Session",
		SID:        "SID",
		ID:         "ID",
		Username:   "Username",
		SigningKey: []byte("secret"),
	}

	tc := token.Create()

	t.Run("GetUserSessions", func(t *testing.T) {

		c, rec := test.MakeRequest("GET", "", false, tc)
		tokenParsed, _ := j.ParseJWT(tc, []byte("secret"))
		c.Set("user", tokenParsed)

		if assert.NoError(t, Sessions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var res paginate.Paginate
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.IsType(t, paginate.Paginate{}, res)
		}
	})

	sessionAfterTest()
}
