package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/Gommunity/GoWithWith/models"
	"github.com/Gommunity/GoWithWith/services/auth"
	"github.com/Gommunity/GoWithWith/services/request"
	"github.com/stretchr/testify/assert"
)

func TestSessions(t *testing.T) {

	ModeliGetUserSessions = func(t1 string, t2, t3 int) (models.Pagination, error) {
		session := &models.Pagination{}
		return *session, nil
	}

	jwt, _ := auth.CreateJWToken("token", "TID", "username", "userID", []byte("secret"))

	t.Run("GetMySessions", func(t *testing.T) {

		c, rec := request.MakeReq("GET", make(url.Values), false, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		if assert.NoError(t, Sessions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON models.Pagination
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Nil(t, errJSON.Data)
			assert.NotNil(t, errJSON.Pages)
			assert.NotNil(t, errJSON.Items)
		}
	})
}
