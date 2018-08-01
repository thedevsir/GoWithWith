package user

import (
	"net/http"
	"os"

	mAuthAttempt "github.com/Gommunity/GoWithWith/app/model/authAttempt"
	mSession "github.com/Gommunity/GoWithWith/app/model/session"
	mUser "github.com/Gommunity/GoWithWith/app/model/user"
	"github.com/Gommunity/GoWithWith/services/auth"
	j "github.com/Gommunity/GoWithWith/services/jwt"
	"github.com/Gommunity/GoWithWith/services/mail"
	"github.com/Gommunity/GoWithWith/services/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	_SendVerficationMail = mail.SendVerficationMail
	_SendResetMail       = mail.SendResetMail
)

type (
	SignupRoute struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=50"`
		Email    string `json:"email" validate:"required,email"`
	}
	ResendRoute struct {
		Email string `json:"email" validate:"required,email"`
	}
	VerificationRoute struct {
		Token string `json:"token" validate:"required"`
	}
	SigninRoute struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	ForgotRoute struct {
		Email string `json:"email" validate:"required,email"`
	}
	ResetRoute struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	PasswordRoute struct {
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
)

// Signup godoc
// @Summary Create an account
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param email body string true "Email"
// @Success 201 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 409 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signup [post]
func Signup(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(SignupRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	userModel := &mUser.User{
		Username: params.Username,
		Password: params.Password,
		Email:    params.Email,
	}

	if _, err = userModel.CheckUsername(userModel.Username); err != nil {
		return r.JSON(http.StatusConflict, err.Error())
	}

	if _, err = userModel.CheckEmail(userModel.Email); err != nil {
		return r.JSON(http.StatusConflict, err.Error())
	}

	userModel.CreateUser()

	token := mail.MakeTokenForEmails("verify", params.Username, params.Email, []byte(os.Getenv("JWTSigningKey")))
	_SendVerficationMail(params.Username, params.Email, token)

	return r.JSON(http.StatusCreated, "Success")
}

// Resend godoc
// @Summary Resend email verfication
// @Tags user
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 404 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signup/resend [post]
func Resend(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(ResendRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	userModel := &mUser.User{
		Email: params.Email,
	}

	var user mUser.User
	if user, err = userModel.CheckEmail(userModel.Email); err != nil {

		if user.VerifyEmail != true {
			return r.JSON(http.StatusBadRequest, "Your account has already been verified")
		}

		token := mail.MakeTokenForEmails("verify", user.Username, params.Email, []byte(os.Getenv("JWTSigningKey")))
		_SendVerficationMail(user.Username, params.Email, token)
		return r.JSON(http.StatusOK, "Success")
	}

	return r.JSON(http.StatusNotFound, "Your email address was not found")
}

// Verification godoc
// @Summary Activate user account
// @Tags user
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Success 200 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signup/verification [post]
func Verification(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(VerificationRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	var data *jwt.Token
	if data, err = j.ParseJWT(params.Token, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	claims := data.Claims.(jwt.MapClaims)
	if claims["Action"].(string) != "verify" {
		return r.JSON(http.StatusBadRequest, "Wrong Action")
	}

	user := &mUser.User{
		Email: claims["Email"].(string),
	}
	user.Activation()

	return r.JSON(http.StatusOK, "Success")
}

// Signin godoc
// @Summary User signin
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 201 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 404 {object} response.Message
// @Failure 429 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signin [post]
func Signin(c echo.Context) (err error) {

	ip := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")

	r := response.Composer{c}
	params := new(SigninRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	userModel := &mUser.User{
		Username: params.Username,
		Password: params.Password,
	}

	userAttempt := &mAuthAttempt.AuthAttempt{
		IP:       ip,
		Username: userModel.Username,
	}

	userAbuse := &mAuthAttempt.AbuseDetected{
		MaxIP:            os.Getenv("AbuseDetectedForIp"),
		MaxIPAndUsername: os.Getenv("AbuseDetectedForIpUsername"),
	}

	if err = userAbuse.Check(ip, userModel.Username); err != nil {
		return r.JSON(http.StatusTooManyRequests, err.Error())
	}

	var user mUser.User
	if user, err = userModel.FindUserByCredentials(); err != nil {
		userAttempt.Create()
		return r.JSON(http.StatusNotFound, err.Error())
	}

	UserID := user.GetId().Hex()

	session := &mSession.Session{
		IP:        ip,
		UserID:    UserID,
		UserAgent: userAgent,
	}
	SID, uuid := session.Create()

	token := &auth.JWTUserToken{
		Session:    uuid,
		SID:        SID,
		ID:         UserID,
		Username:   userModel.Username,
		SigningKey: []byte(os.Getenv("JWTSigningKey")),
	}

	c.Response().Header().Set(echo.HeaderAuthorization, token.Create())
	return r.JSON(http.StatusCreated, "Success")
}

// Forgot godoc
// @Summary Forgot password
// @Tags user
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 404 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signin/forgot [post]
func Forgot(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(ForgotRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	userModel := &mUser.User{
		Email: params.Email,
	}

	var user mUser.User
	if user, err = userModel.CheckEmail(userModel.Email); err != nil {

		token := mail.MakeTokenForEmails("reset", user.Username, params.Email, []byte(os.Getenv("JWTSigningKey")))
		_SendResetMail(user.Username, params.Email, token)
		return r.JSON(http.StatusOK, "Success")
	}

	return r.JSON(http.StatusNotFound, "Your email address was not found")
}

// Reset godoc
// @Summary Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/signin/reset [post]
func Reset(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(ResetRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	var data *jwt.Token
	if data, err = j.ParseJWT(params.Token, []byte(os.Getenv("JWTSigningKey"))); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	claims := data.Claims.(jwt.MapClaims)
	if claims["Action"].(string) != "reset" {
		return r.JSON(http.StatusBadRequest, "Wrong Action")
	}

	user := &mUser.User{
		Email:    claims["Email"].(string),
		Password: params.Password,
	}
	user.ChangePassword()

	return r.JSON(http.StatusOK, "Success")
}

// Password godoc
// @Summary Change password
// @Tags user
// @Accept json
// @Produce json
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/auth/password [put]
func Password(c echo.Context) (err error) {

	r := response.Composer{c}
	params := new(PasswordRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	data := c.Get("user").(*jwt.Token)
	claims := data.Claims.(jwt.MapClaims)

	user := &mUser.User{
		Username: claims["userId"].(string),
		Password: params.Password,
	}
	user.ChangePassword()

	return r.JSON(http.StatusOK, "Success")
}
