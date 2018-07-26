package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Gommunity/GoWithWith/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	form := make(url.Values)
	form.Set("username", "irani")
	form.Set("password", "12345")

	ModeliAbuseDetected = models.AbuseDetected{
		MaxIP:            "0",
		MaxIPAndUsername: "0",
	}
	ModeliAbuseDetectedCheck = func(t1, t2 string) error {
		return nil
	}
	ModeliFindUserByCredentials = func(t1, t2 string) (models.User, error) {
		user := &models.User{}
		return *user, nil
	}
	ModeliCreateSession = func(t1, t2, t3 string) (string, string, error) {
		return "", "", nil
	}

	t.Run("Success", func(t *testing.T) {

		c, rec := MakeReq("POST", form, true, "")

		if assert.NoError(t, Login(c)) {
			var res Authorization
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.NotNil(t, res.Authorization)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Failure", func(t *testing.T) {

		form.Set("password", "")

		c, rec := MakeReq("POST", form, true, "")

		if assert.NoError(t, Login(c)) {
			var errJSON JoiError
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "cannot be blank", errJSON.Message["Password"])
		}
	})
}

func TestForgot(t *testing.T) {

	form := make(url.Values)
	form.Set("email", "p30search@gmail.com")

	ModeliCheckEmail = func(t1 string) (models.User, error) {
		return models.User{}, errors.New("")
	}
	ModeliMakeEmailToken = func(t1, t2, t3 string, t4 []byte) (string, error) {
		return "", nil
	}
	ModeliSendResetMail = func(t1, t2, t3 string) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {

		c, rec := MakeReq("POST", form, true, "")

		if assert.NoError(t, Forgot(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}

func TestReset(t *testing.T) {

	token, _ := ModeliMakeEmailToken("reset", "username", "email", []byte("secret"))

	form := make(url.Values)
	form.Set("token", "fakeToken")
	form.Set("password", "blahblahblah")

	ModeliParseJWT = func(t1 string, t2 []byte) (*jwt.Token, error) {
		return ParseJWT(token, []byte("secret"))
	}
	ModeliChangeUserPassword = func(t1, t2 string) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {

		c, rec := MakeReq("POST", form, true, "")

		if assert.NoError(t, Reset(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}
