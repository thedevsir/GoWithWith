package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	helpers "../../helpers"
	models "../../models"
	structs "../structs"
	"github.com/labstack/echo"
)

// Login godoc
// @Summary User login
// @Description Login and get jwt session
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} structs.Authorization
// @Failure 400 {object} helpers.JoiError
// @Router /user/login [post]
func Login(c echo.Context) (err error) {

	var user models.User
	var tokenJWT, session, SID, userID string

	ip := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")

	Logi := structs.LoginStruct{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if err = c.Bind(new(structs.LoginStruct)); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = Logi.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if err = ModeliAbuseDetectedCheck(ip, Logi.Username); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	if user, err = ModeliFindUserByCredentials(Logi.Username, Logi.Password); err != nil {
		models.AttemptCreate(ip, Logi.Username)
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	userID = user.GetId().Hex()

	if SID, session, err = ModeliCreateSession(userID, ip, userAgent); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	if tokenJWT, err = ModeliCreateJWToken(session, SID, user.Username, userID, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	auth := &structs.Authorization{
		Authorization: tokenJWT,
	}

	return c.JSON(http.StatusOK, auth)
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
// @Router /user/forgot [post]
func Forgot(c echo.Context) (err error) {

	var token string
	var user models.User

	forgo := structs.ForgotStruct{
		Email: c.FormValue("email"),
	}

	if err = c.Bind(new(structs.ForgotStruct)); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = forgo.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if user, err = ModeliCheckEmail(forgo.Email); err != nil {

		token, err = ModeliMakeEmailToken(user.Username, user.Email, []byte(os.Getenv("JWTSigningKey")))

		if err = ModeliSendResetMail(user.Username, user.Email, token); err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
		}

		return c.JSON(http.StatusOK, helpers.SayOk("Success."))
	}

	return c.JSON(http.StatusOK, helpers.ThrowString(errors.New("Email not found")))
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
// @Router /user/reset [post]
func Reset(c echo.Context) (err error) {

	var data *jwt.Token

	rese := structs.ResetStruct{
		Token:    c.FormValue("token"),
		Password: c.FormValue("password"),
	}

	if err = c.Bind(new(structs.ForgotStruct)); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = rese.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if data, err = ModeliParseJWT(rese.Token, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	// Get data from token
	claims := data.Claims.(jwt.MapClaims)

	if err = ModeliChangeUserPassword(claims["Username"].(string), rese.Password); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	return c.JSON(http.StatusOK, helpers.SayOk("Success."))
}
