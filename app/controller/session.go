package controller

import (
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/Gommunity/GoWithWith/services/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Sessions godoc
// @Summary Get user sessions
// @Description Get all user sessions
// @Tags users
// @Produce  json
// @Param page query number false "Page"
// @Param limit query number false "Limit"
// @Security ApiKeyAuth
// @Success 200 {string} repository.Pagination
// @Failure 400 {object} services.JoiString
// @Router /user/auth/sessions [get]
func Sessions(c echo.Context) error {

	var err error
	var sessions repository.Pagination

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["userId"].(string)

	page, limit := paginate.Settings(c)

	if sessions, err = repository.GetUserSessions(UserID, page, limit); err != nil {
		return response.Error(err.Error(), 1005)
	}
	return response.Data(c, sessions)
}
