package session

import (
	"net/http"

	mSession "github.com/Gommunity/GoWithWith/app/model/session"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/Gommunity/GoWithWith/services/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	LogoutRoute struct {
		ID string `json:"id" validate:"required"`
	}
)

// Sessions godoc
// @Summary Get user sessions
// @Tags session
// @Accept json
// @Produce json
// @Param page query number false "Page"
// @Param limit query number false "Limit"
// @Security ApiKeyAuth
// @Success 200 {string} response.Message
// @Failure 404 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/v1/auth/sessions [get]
func Sessions(c echo.Context) error {

	var err error
	var sessions paginate.Paginate
	r := response.Composer{c}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)

	page, limit := paginate.HandleQueries(c)
	if sessions, err = mSession.GetUserSessions(userID, page, limit); err != nil {
		return r.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, sessions)
}

// Logout godoc
// @Summary Delete session
// @Tags session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id body string false "Session ID"
// @Success 200 {string} response.Message
// @Failure 404 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /user/auth/logout [delete]
func Logout(c echo.Context) error {

	var session mSession.Session

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var userID, SID = claims["userId"].(string), claims["sid"].(string)

	r := response.Composer{c}
	params := new(LogoutRoute)

	if err := c.Bind(params); err != nil {
		return r.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err == nil {
		if session, err = mSession.SessionFindByID(params.ID); err != nil {
			return r.JSON(http.StatusNotFound, err.Error())
		}
		if session.UserID == userID {
			SID = params.ID
		}
	}

	mSession.DeleteSession(SID)

	return r.JSON(http.StatusOK, "Success")
}
