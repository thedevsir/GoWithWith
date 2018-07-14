package handlers

import (
	"os"

	"github.com/Gommunity/GoWithWith/helpers/response"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo"
)

type SignupStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s SignupStruct) Joi() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required, validation.Length(3, 50), is.Alphanumeric),
		validation.Field(&s.Password, validation.Required, validation.Length(3, 50)),
		validation.Field(&s.Email, validation.Required, is.Email),
	)
}

// Signup godoc
// @Summary Create a account
// @Description Create by multipart/form-data
// @Tags users
// @Accept mpfd
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Success 200 {object} helpers.JoiString
// @Failure 400 {object} helpers.JoiError
// @Router /user/signup [post]
func Signup(c echo.Context) (err error) {

	params := SignupStruct{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	if err = params.Joi(); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if _, err = ModeliCheckUsername(params.Username); err != nil {
		return response.Error(err.Error(), 1002)
	}

	if _, err = ModeliCheckEmail(params.Email); err != nil {
		return response.Error(err.Error(), 1003)
	}

	ModeliCreateUser(params.Username, params.Password, params.Email)

	token := ModeliMakeEmailToken("verify", params.Username, params.Email, []byte(os.Getenv("JWTSigningKey")))
	ModeliSendVerficationMail(params.Username, params.Email, token)

	return response.Created(c, "Created Successfully")
}

type ResendEmailStruct struct {
	Email string `json:"email"`
}

func (r ResendEmailStruct) Joi() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}

// ResendEmail godoc
// @Summary Resend verfication email
// @Description Create by multipart/form-data
// @Tags users
// @Accept mpfd
// @Produce json
// @Param email formData string true "Email"
// @Success 200 {object} helpers.JoiString
// @Failure 400 {object} helpers.JoiError
// @Router /user/signup/resend-email [post]
func ResendEmail(c echo.Context) (err error) {

	params := ResendEmailStruct{
		Email: c.FormValue("email"),
	}

	if err = params.Joi(); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if _, err = ModeliCheckEmailVerify(params.Email); err != nil {
		return response.Error(err.Error(), 1004)
	}

	return response.Ok(c, "Success")
}
