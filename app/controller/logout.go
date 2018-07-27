package controller

import (
	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/services/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Logout godoc
// @Summary Logout user
// @Description Delete current session or delete special session with id
// @Tags users
// @Produce  json
// @Security ApiKeyAuth
// @Param id formData string false "Session ID"
// @Success 200 {string} services.JoiString
// @Failure 400 {object} services.JoiError
// @Router /user/auth/logout [delete]
func Logout(c echo.Context) error {

	var session repository.Session

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["userId"].(string)
	SID := claims["sid"].(string)

	params := new(model.LogoutStruct)

	if err := c.Bind(params); err != nil {
		return response.Error(err.Error(), 1000)
	}

	if err := c.Validate(params); err == nil {
		if session, err = repository.SessionFindByID(params.ID); err != nil {
			return response.Error(err.Error(), 1006)
		}
		if session.UserID == UserID {
			SID = params.ID
		}
	}

	repository.DeleteSession(SID)
	return response.Ok(c, "Successfully Signout")
}
