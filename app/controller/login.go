package controller

import (
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/config/auth"
	"github.com/Gommunity/GoWithWith/services/mail"
	"github.com/Gommunity/GoWithWith/services/response"
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
// @Success 200 {object} Authorization
// @Failure 400 {object} services.JoiError
// @Router /user/login [post]
func Login(c echo.Context) (err error) {

	var user repository.User
	var tokenJWT, session, SID, userID string

	ip := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")

	params := new(model.LoginStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if err = auth.AbuseDetected.Check(ip, params.Username); err != nil {
		return response.Error(err.Error(), 1007)
	}

	if user, err = repository.FindUserByCredentials(params.Username, params.Password); err != nil {
		repository.AttemptCreate(ip, params.Username)
		return response.Error(err.Error(), 1008)
	}

	userID = user.GetId().Hex()

	SID, session = repository.CreateSession(userID, ip, userAgent)

	tokenJWT = auth.CreateJWToken(session, SID, user.Username, userID, []byte(os.Getenv("JWTSigningKey")))

	auth := &model.Authorization{
		Authorization: tokenJWT,
	}
	return response.Data(c, auth)
}

// Forgot godoc
// @Summary Forgot password
// @Description Request for reset password
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param email formData string true "Email"
// @Success 200 {object} services.SayOk
// @Failure 400 {object} services.JoiError
// @Router /user/login/forgot [post]
func Forgot(c echo.Context) (err error) {

	var token string
	var user repository.User

	params := new(model.ForgotStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if user, err = repository.CheckEmail(params.Email); err != nil {

		token = mail.MakeEmailToken("reset", user.Username, user.Email, []byte(os.Getenv("JWTSigningKey")))

		mail.SendResetMail(user.Username, user.Email, token)

		return response.Ok(c, "Success")
	}
	return response.Error("Email not found", 1009)
}

// Reset godoc
// @Summary Reset password
// @Description Change account password
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param token formData string true "Token"
// @Param password formData string true "Password"
// @Success 200 {object} services.SayOk
// @Failure 400 {object} services.JoiError
// @Router /user/login/reset [post]
func Reset(c echo.Context) (err error) {

	var data *jwt.Token

	params := new(model.ResetStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if data, err = auth.ParseJWT(params.Token, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return response.Error(err.Error(), 1010)
	}

	claims := data.Claims.(jwt.MapClaims)

	if claims["Action"].(string) != "reset" {
		return response.Error("wrong action type", 1011)
	}

	repository.ChangeUserPassword(claims["Username"].(string), params.Password)

	return response.Ok(c, "Success")
}
