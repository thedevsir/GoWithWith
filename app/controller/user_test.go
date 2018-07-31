package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/mail"
	"github.com/Gommunity/GoWithWith/services/response"
	"github.com/Gommunity/GoWithWith/services/test"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/zebresel-com/mongodm"
)

var userCollection *mongodm.Model

func userBeforeTest() {

	utility.LoadEnvironmentVariables("../../.env")

	db := test.DBComposer("../../resource/locals/locals.json")
	db.Shoot(map[string]mongodm.IDocumentBase{
		"authAttempts": &model.AuthAttempt{},
		"sessions":     &model.Session{},
		"users":        &model.User{},
	})

	userCollection = database.Connection.Model(model.UserCollection)
	userCollection.RemoveAll(nil)
}

func userAfterTest() {
	userCollection.RemoveAll(nil)
}

func TestSignup(t *testing.T) {

	userBeforeTest()

	// Mock
	_SendVerficationMail = func(username, email, token string) {}

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, ``, false, "")

		if assert.NoError(t, Signup(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"username":"ir","password":"12345","email":"fakeMail"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signup(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("DuplicateUsername", func(t *testing.T) {

		JSONData := `{"username":"Amir","password":"12345678","email":"live@forumX.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signup(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusConflict, rec.Code)
			assert.Equal(t, "The username has already been used", res.Message)
		}
	})

	t.Run("DuplicateEmail", func(t *testing.T) {

		JSONData := `{"username":"Amir2","password":"12345678","email":"live@forumX.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signup(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusConflict, rec.Code)
			assert.Equal(t, "The email address has already been used", res.Message)
		}
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := `{"username":"irani","password":"12345678","email":"fakeMail@gmail.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signup(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	userAfterTest()
}

func TestResend(t *testing.T) {

	userBeforeTest()

	// Mock
	_SendVerficationMail = func(username, email, token string) {}

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, ``, false, "")

		if assert.NoError(t, Resend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Resend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("EmailNotFound", func(t *testing.T) {

		JSONData := `{"email":"fake@gmail.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Resend(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, "Your email address was not found", res.Message)
		}
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := `{"email":"live@forumX.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Resend(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	t.Run("VerfiedBefore", func(t *testing.T) {

		userModel.Activation()

		JSONData := `{"email":"live@forumX.com"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Resend(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Your account has already been verified", res.Message)
		}
	})

	userAfterTest()
}

func TestVerification(t *testing.T) {

	userBeforeTest()

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	token0 := mail.MakeTokenForEmails("verify", "username", "live@forumX.com", []byte(os.Getenv("JWTSigningKey")))
	token1 := mail.MakeTokenForEmails("reset", "username", "live@forumX.com", []byte(os.Getenv("JWTSigningKey")))

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, ``, false, "")

		if assert.NoError(t, Verification(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Verification(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("WrongAction", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"token":"%s"}`, token1)
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Verification(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Wrong Action", res.Message)
		}
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"token":"%s"}`, token0)
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Verification(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	userAfterTest()
}

func TestSignin(t *testing.T) {

	userBeforeTest()

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, ``, false, "")

		if assert.NoError(t, Signin(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"username":"ir"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signup(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {

		JSONData := `{"username":"ir","password":"12345"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signin(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("LoginSuccessfully", func(t *testing.T) {

		JSONData := `{"username":"Amir","password":"12345678"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Signin(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	userAfterTest()
}

func TestForgot(t *testing.T) {

	userBeforeTest()

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	// Mock
	_SendResetMail = func(username, email, token string) {}

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, ``, false, "")

		if assert.NoError(t, Forgot(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"key":"value"}`
		c, rec := test.MakeRequest(echo.POST, JSONData, true, "")

		if assert.NoError(t, Forgot(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("EmailNotFound", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, `{"email":"fake@fake.io"}`, true, "")

		if assert.NoError(t, Forgot(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, "Your email address was not found", res.Message)
		}
	})

	t.Run("Success", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.POST, `{"email":"live@forumX.com"}`, true, "")

		if assert.NoError(t, Forgot(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	userAfterTest()
}

func TestReset(t *testing.T) {

	userBeforeTest()

	// Create user
	userModel := &model.User{
		Username: "Amir",
		Password: "12345678",
		Email:    "live@forumX.com",
	}
	userModel.CreateUser()

	token0 := mail.MakeTokenForEmails("reset", "username", "live@forumX.com", []byte(os.Getenv("JWTSigningKey")))
	token1 := mail.MakeTokenForEmails("verify", "username", "live@forumX.com", []byte(os.Getenv("JWTSigningKey")))

	t.Run("BindErr", func(t *testing.T) {

		c, rec := test.MakeRequest(echo.PUT, ``, false, "")

		if assert.NoError(t, Reset(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ValidateErr", func(t *testing.T) {

		JSONData := `{"password":"12345"}`
		c, rec := test.MakeRequest(echo.PUT, JSONData, true, "")

		if assert.NoError(t, Reset(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("WrongAction", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"password":"12345678","token":"%s"}`, token1)
		c, rec := test.MakeRequest(echo.PUT, JSONData, true, "")

		if assert.NoError(t, Reset(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Wrong Action", res.Message)
		}
	})

	t.Run("Success", func(t *testing.T) {

		JSONData := fmt.Sprintf(`{"password":"12345678","token":"%s"}`, token0)
		c, rec := test.MakeRequest(echo.PUT, JSONData, true, "")

		if assert.NoError(t, Reset(c)) {
			var res response.Message
			json.Unmarshal([]byte(rec.Body.String()), &res)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "Success", res.Message)
		}
	})

	userAfterTest()
}
