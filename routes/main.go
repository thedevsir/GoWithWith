package routes

import (
	"os"

	cs "github.com/Gommunity/GoWithWith/app/controller/session"
	cu "github.com/Gommunity/GoWithWith/app/controller/user"
	"github.com/Gommunity/GoWithWith/middleware/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
)

func Composer() *echo.Echo {

	Route := echo.New()
	endpoints := Route.Group("/endpoint")
	{
		User := endpoints.Group("/user/v1")
		{
			User.POST("/signup", cu.Signup)
			User.POST("/signup/resend", cu.Resend)
			User.POST("/signup/verification", cu.Verification)
			User.POST("/signin", cu.Signin)
			User.POST("/signin/forgot", cu.Forgot)
			User.PUT("/signin/reset", cu.Reset)
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(os.Getenv("JWTSigningKey"))))
				Auth.Use(auth.Middleware)
				Auth.PUT("/password", cu.Password)
				Auth.GET("/sessions", cs.Sessions)
				Auth.DELETE("/logout", cs.Logout)
			}
		}
	}
	Route.GET("/swagger/*", echoSwagger.WrapHandler)
	return Route
}
