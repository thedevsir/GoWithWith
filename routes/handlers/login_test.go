package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	models "../../models"
	structs "../structs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	// From
	form := make(url.Values)
	form.Set("username", "irani")
	form.Set("password", "12345")

	// Mock
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

		// Setup
		c, rec := MakeReq("POST", form, true, "")

		// Assertions
		if assert.NoError(t, Login(c)) {
			var res structs.Authorization
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.NotNil(t, res.Authorization)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Failure", func(t *testing.T) {

		form.Set("password", "")

		// Setup
		c, rec := MakeReq("POST", form, true, "")

		// Assertions
		if assert.NoError(t, Login(c)) {
			var errJSON JoiError
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "cannot be blank", errJSON.Message["Password"])
		}
	})
}

func TestForgot(t *testing.T) {

	// From
	form := make(url.Values)
	form.Set("email", "p30search@gmail.com")

	// Mock
	ModeliCheckEmail = func(t1 string) (models.User, error) {
		return models.User{}, errors.New("")
	}
	ModeliMakeEmailToken = func(t1, t2 string, t3 []byte) (string, error) {
		return "", nil
	}
	ModeliSendResetMail = func(t1, t2, t3 string) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {

		// Setup
		c, rec := MakeReq("POST", form, true, "")

		// Assertions
		if assert.NoError(t, Forgot(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}

func TestReset(t *testing.T) {

	token, _ := ModeliMakeEmailToken("username", "email", []byte("secret"))

	// From
	form := make(url.Values)
	form.Set("token", "fakeToken")
	form.Set("password", "blahblahblah")

	// Mock
	ModeliParseJWT = func(t1 string, t2 []byte) (*jwt.Token, error) {
		return ParseJWT(token, []byte("secret"))
	}
	ModeliChangeUserPassword = func(t1, t2 string) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {

		// Setup
		c, rec := MakeReq("POST", form, true, "")

		// Assertions
		if assert.NoError(t, Reset(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var errJSON JoiString
			json.Unmarshal([]byte(rec.Body.String()), &errJSON)
			assert.Equal(t, "Success.", errJSON.Message)
		}
	})
}
