package auth

import (
	"net/http"

	helpers "../helpers"
	models "../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// AuthMiddleware ...
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		TID := claims["tid"].(string)
		token := claims["token"].(string)

		if err := models.SessionFindByCredentials(token, TID); err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
		}

		models.UpdateLastActive(TID)

		return next(c)
	}
}
