package controller

import (
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/helpers/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type LogoutStruct struct {
	ID string `json:"id" validate:"required"`
}

// Logout godoc
// @Summary Logout user
// @Description Delete current session or delete special session with id
// @Tags users
// @Produce  json
// @Security ApiKeyAuth
// @Param id formData string false "Session ID"
// @Success 200 {string} helpers.JoiString
// @Failure 400 {object} helpers.JoiError
// @Router /user/auth/logout [delete]
func Logout(c echo.Context) error {

	var session repository.Session

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["userId"].(string)
	SID := claims["sid"].(string)

	params := new(LogoutStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err == nil {
		if session, err = ModeliSessionFindByID(params.ID); err != nil {
			return response.Error(err.Error(), 1006)
		}
		if session.UserID == UserID {
			SID = params.ID
		}
	}

	ModeliDeleteSession(SID)

	return response.Ok(c, "Successfully Signout")
}
