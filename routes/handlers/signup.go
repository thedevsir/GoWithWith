package handlers

import (
	"net/http"

	helpers "../../helpers"
	structs "../structs"
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
// @Success 200 {object} helpers.JoiString
// @Failure 400 {object} helpers.JoiError
// @Router /user/signup [post]
func Signup(c echo.Context) (err error) {

	Sign := structs.SignupStruct{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	if err = c.Bind(new(structs.SignupStruct)); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = Sign.Joi(); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if _, err = ModeliCheckUsername(Sign.Username); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if _, err = ModeliCheckEmail(Sign.Email); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	if err = ModeliCreateUser(Sign.Username, Sign.Password, Sign.Email); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Throw(err))
	}

	return c.JSON(http.StatusOK, helpers.SayOk("Success."))
}
