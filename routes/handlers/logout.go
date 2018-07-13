package handlers

import (
	"net/http"

	"github.com/Gommunity/GoWithWith/helpers"
	"github.com/Gommunity/GoWithWith/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

type LogoutStruct struct {
	ID string `form:"id"`
}

func (l LogoutStruct) Joi() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.ID, validation.Required),
	)
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

	var session models.Session

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["userId"].(string)
	SID := claims["sid"].(string)

	params := LogoutStruct{
		ID: c.FormValue("id"),
	}

	err := params.Joi()

	if err == nil {
		if session, err = ModeliSessionFindByID(params.ID); err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
		}
		if session.UserID == UserID {
			SID = params.ID
		}
	}

	ModeliDeleteSession(SID)

	return c.JSON(http.StatusOK, helpers.SayOk("Success."))
}
