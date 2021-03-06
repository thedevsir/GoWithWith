package auth

import (
	"net/http"

	ms "github.com/Gommunity/GoWithWith/app/model/session"
	"github.com/Gommunity/GoWithWith/services/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		r := response.Composer{c}
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		var SID, session = claims["sid"].(string), claims["session"].(string)

		if err := ms.SessionFindByCredentials(session, SID); err != nil {
			return r.JSON(http.StatusBadRequest, err.Error())
		}

		ms.SessionUpdateLastActivity(SID)

		return next(c)
	}
}
