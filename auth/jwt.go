package auth

import (
	"net/http"

	"github.com/Gommunity/GoWithWith/helpers"
	"github.com/Gommunity/GoWithWith/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		var SID, session = claims["sid"].(string), claims["session"].(string)

		if err := models.SessionFindByCredentials(session, SID); err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ThrowString(err))
		}

		models.UpdateLastActive(SID)

		return next(c)
	}
}
