package routes

import (
	"os"

	c "github.com/Gommunity/GoWithWith/app/controller"
	"github.com/Gommunity/GoWithWith/middleware/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
)

func Initial() *echo.Echo {

	Route := echo.New()
	endpoints := Route.Group("/endpoint")
	{
		User := endpoints.Group("/user/v1")
		{
			User.POST("/signup", c.Signup)
			User.POST("/signup/resend", c.Resend)
			User.POST("/signup/verification", c.Verification)
			User.POST("/signin", c.Signin)
			User.POST("/signin/forgot", c.Forgot)
			User.PUT("/signin/reset", c.Reset)
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(os.Getenv("JWTSigningKey"))))
				Auth.Use(auth.Middleware)
				Auth.GET("/sessions", c.Sessions)
				Auth.DELETE("/logout", c.Logout)
			}
		}
	}
	Route.GET("/swagger/*", echoSwagger.WrapHandler)
	return Route
}
