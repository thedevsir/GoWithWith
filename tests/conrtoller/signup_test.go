package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/Gommunity/GoWithWith/models"
	"github.com/Gommunity/GoWithWith/services/request"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {

	form := make(url.Values)
	form.Set("password", "12345")

	ModeliCheckUsername = func(t1 string) (models.User, error) {
		return models.User{}, nil
	}
	ModeliCheckEmail = func(t1 string) (models.User, error) {
		return models.User{}, nil
	}
	ModeliCreateUser = func(t1, t2, t3 string) error {
		return nil
	}
	ModeliMakeEmailToken = func(t1, t2, t3 string, t4 []byte) (string, error) {
		return "", nil
	}
	ModeliSendVerficationMail = func(t1, t2, t3 string) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {

		form.Set("username", "irani")
		form.Set("email", "freshmanlimited@gmail.com")

		c, rec := request.MakeReq("POST", form, true, "")

		if assert.NoError(t, Signup(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})

	t.Run("Failure", func(t *testing.T) {

		form.Set("username", "ir") // min: 3
		form.Set("email", "wrongEmail")

		c, rec := request.MakeReq("POST", form, true, "")

		if assert.NoError(t, Signup(c)) {
			var errJSON JoiError
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "must be a valid email address", errJSON.Message["Email"])
			assert.Equal(t, "the length must be between 3 and 50", errJSON.Message["Username"])
		}
	})
}
