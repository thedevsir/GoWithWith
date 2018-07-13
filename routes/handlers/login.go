package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/Gommunity/GoWithWith/helpers"
	"github.com/Gommunity/GoWithWith/models"
	"github.com/labstack/echo"
)

type LoginStruct struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type Authorization struct {
	Authorization string `json:"authorization"`
}

func (l LoginStruct) Joi() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Username, validation.Required),
		validation.Field(&l.Password, validation.Required),
	)
}

// Login godoc
// @Summary User login
// @Description Login and get jwt session
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} Authorization
// @Failure 400 {object} helpers.JoiError
// @Router /user/login [post]
func Login(c echo.Context) (err error) {

	var user models.User
	var tokenJWT, session, SID, userID string

	ip := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")

	params := LoginStruct{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if err = params.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if err = ModeliAbuseDetectedCheck(ip, params.Username); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	if user, err = ModeliFindUserByCredentials(params.Username, params.Password); err != nil {
		models.AttemptCreate(ip, params.Username)
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	userID = user.GetId().Hex()

	SID, session = ModeliCreateSession(userID, ip, userAgent)

	tokenJWT = ModeliCreateJWToken(session, SID, user.Username, userID, []byte(os.Getenv("JWTSigningKey")))

	auth := &Authorization{
		Authorization: tokenJWT,
	}

	return c.JSON(http.StatusOK, auth)
}

type ForgotStruct struct {
	Email string `form:"email"`
}

func (f ForgotStruct) Joi() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Email, validation.Required, is.Email),
	)
}

// Forgot godoc
// @Summary Forgot password
// @Description Request for reset password
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param email formData string true "Email"
// @Success 200 {object} helpers.SayOk
// @Failure 400 {object} helpers.JoiError
// @Router /user/login/forgot [post]
func Forgot(c echo.Context) (err error) {

	var token string
	var user models.User

	params := ForgotStruct{
		Email: c.FormValue("email"),
	}

	if err = params.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if user, err = ModeliCheckEmail(params.Email); err != nil {

		token = ModeliMakeEmailToken("reset", user.Username, user.Email, []byte(os.Getenv("JWTSigningKey")))

		ModeliSendResetMail(user.Username, user.Email, token)

		return c.JSON(http.StatusOK, helpers.SayOk("Success."))
	}

	return c.JSON(http.StatusOK, helpers.ThrowString(errors.New("Email not found")))
}

type ResetStruct struct {
	Token    string `form:"token"`
	Password string `form:"password"`
}

func (r ResetStruct) Joi() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(3, 50)),
	)
}

// Reset godoc
// @Summary Reset password
// @Description Change account password
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param token formData string true "Token"
// @Param password formData string true "Password"
// @Success 200 {object} helpers.SayOk
// @Failure 400 {object} helpers.JoiError
// @Router /user/login/reset [post]
func Reset(c echo.Context) (err error) {

	var data *jwt.Token

	params := ResetStruct{
		Token:    c.FormValue("token"),
		Password: c.FormValue("password"),
	}

	if err = params.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if data, err = ModeliParseJWT(params.Token, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	claims := data.Claims.(jwt.MapClaims)

	if claims["Action"].(string) != "reset" {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(errors.New("wrong action type")))
	}

	ModeliChangeUserPassword(claims["Username"].(string), params.Password)

	return c.JSON(http.StatusOK, helpers.SayOk("Success."))
}
