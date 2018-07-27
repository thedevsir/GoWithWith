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

func TestLogout(t *testing.T) {

	form := make(url.Values)
	form.Set("id", "5b169a9cb8c97417fa5e0c28") // fakeID

	ModeliSessionFindByID = func(t1 string) (models.Session, error) {
		session := &models.Session{
			UserID: "UserID",
		}
		return *session, nil
	}
	ModeliDeleteSession = func(t1 string) error {
		return nil
	}

	jwt, _ := auth.CreateJWToken("session", "SID", "username", "userID", []byte("secret"))

	t.Run("DeleteCurrentSession", func(t *testing.T) {

		c, rec := request.MakeReq("DELETE", form, false, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		if assert.NoError(t, Logout(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})

	t.Run("DeleteByID", func(t *testing.T) { // Cant detect id in postData

		c, rec := request.MakeReq("DELETE", form, true, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		if assert.NoError(t, Logout(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}
