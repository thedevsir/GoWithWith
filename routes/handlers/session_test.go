package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	models "../../models"
	"github.com/stretchr/testify/assert"
)

func TestSessions(t *testing.T) {

	// Mock
	ModeliGetUserSessions = func(t1 string, t2, t3 int) (models.Pagination, error) {
		session := &models.Pagination{}
		return *session, nil
	}

	// GenerateJWT
	jwt, _ := models.CreateJWToken("token", "TID", "username", "userID", []byte("secret"))

	t.Run("GetMySessions", func(t *testing.T) {

		// Setup
		c, rec := MakeReq("GET", make(url.Values), false, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		// Assertions
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
