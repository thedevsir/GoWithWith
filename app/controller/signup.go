package controller

import (
	"os"

	"github.com/Gommunity/GoWithWith/helpers/response"
	"github.com/labstack/echo"
)

type SignupStruct struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type ResendEmailStruct struct {
	Email string `json:"email" validate:"required,email"`
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

	params := new(SignupStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
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

// func Verify(c echo.Context) (err error) {

// 	// params := new()
// }

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

	params := new(ResendEmailStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if _, err = ModeliCheckEmailVerify(params.Email); err != nil {
		return response.Error(err.Error(), 1004)
	}

	return response.Ok(c, "Success")
}
