package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	models "../../models"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {

	// From
	form := make(url.Values)
	form.Set("id", "5b169a9cb8c97417fa5e0c28") // fakeID

	// Mock
	ModeliSessionFindByID = func(t1 string) (models.Session, error) {
		session := &models.Session{
			UserID: "UserID",
		}
		return *session, nil
	}
	ModeliDeleteSession = func(t1 string) error {
		return nil
	}

	// GenerateJWT
	jwt, _ := models.CreateJWToken("session", "SID", "username", "userID", []byte("secret"))

	t.Run("DeleteCurrentSession", func(t *testing.T) {

		// Setup
		c, rec := MakeReq("DELETE", form, false, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		// Assertions
		if assert.NoError(t, Logout(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})

	t.Run("DeleteByID", func(t *testing.T) { // Cant detect id in postData

		// Setup
		c, rec := MakeReq("DELETE", form, true, jwt)
		token, _ := ParseJWT(jwt, []byte("secret"))
		c.Set("user", token)

		// Assertions
		if assert.NoError(t, Logout(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}