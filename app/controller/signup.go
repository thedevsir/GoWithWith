package controller

import (
	"os"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/services/mail"
	"github.com/Gommunity/GoWithWith/services/response"
	"github.com/labstack/echo"
)

// Signup godoc
// @Summary Create a account
// @Description Create by multipart/form-data
// @Tags users
// @Accept mpfd
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Success 200 {object} services.JoiString
// @Failure 400 {object} services.JoiError
// @Router /user/signup [post]
func Signup(c echo.Context) (err error) {

	params := new(model.SignupStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if _, err = repository.CheckUsername(params.Username); err != nil {
		return response.Error(err.Error(), 1002)
	}

	if _, err = repository.CheckEmail(params.Email); err != nil {
		return response.Error(err.Error(), 1003)
	}

	repository.CreateUser(params.Username, params.Password, params.Email)
	token := mail.MakeEmailToken("verify", params.Username, params.Email, []byte(os.Getenv("JWTSigningKey")))
	mail.SendVerficationMail(params.Username, params.Email, token)
	return response.Created(c, "Created Successfully")
}

// ResendEmail godoc
// @Summary Resend verfication email
// @Description Create by multipart/form-data
// @Tags users
// @Accept mpfd
// @Produce json
// @Param email formData string true "Email"
// @Success 200 {object} services.JoiString
// @Failure 400 {object} services.JoiError
// @Router /user/signup/resend-email [post]
func ResendEmail(c echo.Context) (err error) {

	params := new(model.ResendEmailStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err != nil {
		return response.Error(err.Error(), 1001)
	}

	if _, err = repository.CheckEmailVerify(params.Email); err != nil {
		return response.Error(err.Error(), 1004)
	}

	return response.Ok(c, "Success")
}
