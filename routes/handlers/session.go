package handlers

import (
	"github.com/Gommunity/GoWithWith/helpers/response"
	"github.com/Gommunity/GoWithWith/models"
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
// @Success 200 {string} models.Pagination
// @Failure 400 {object} helpers.JoiString
// @Router /user/auth/sessions [get]
func Sessions(c echo.Context) error {

	var err error
	var sessions models.Pagination

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["userId"].(string)

	page, limit := PaginationSettings(c)

	if sessions, err = ModeliGetUserSessions(UserID, page, limit); err != nil {
		return response.Error(err.Error(), 1005)
	}
	return response.Data(c, sessions)
}
