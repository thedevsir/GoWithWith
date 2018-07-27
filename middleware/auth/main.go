package auth

import (
	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/Gommunity/GoWithWith/services/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		var SID, session = claims["sid"].(string), claims["session"].(string)

		if err := repository.SessionFindByCredentials(session, SID); err != nil {
			return response.Error(err.Error(), 1012)
		}

		repository.UpdateLastActive(SID)

		return next(c)
	}
}
