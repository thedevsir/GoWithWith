package handlers

import (
	"net/http"

	helpers "../../helpers"
	models "../../models"
	structs "../structs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

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

	Logout := structs.LogoutStruct{
		ID: c.FormValue("id"),
	}

	err := Logout.Joi()

	if err == nil {
		if session, err = ModeliSessionFindByID(Logout.ID); err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
		}
		if session.UserID == UserID {
			SID = Logout.ID
		}
	}

	if err = ModeliDeleteSession(SID); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
	}

	return c.JSON(http.StatusOK, helpers.SayOk("Success."))
}
