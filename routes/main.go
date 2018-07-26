package routes

import (
	"os"

	"github.com/Gommunity/GoWithWith/app/controller"
	"github.com/Gommunity/GoWithWith/middleware/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
)

func Initial() *echo.Echo {

	Route := echo.New()

	endpointV1 := Route.Group("/endpoint/v1")
	{
		User := endpointV1.Group("/user")
		{
			User.POST("/signup", controller.Signup)
			// User.POST("/signup/verify", controller.Verify)
			User.POST("/signup/resend-email", controller.ResendEmail)
			User.POST("/login", controller.Login)
			User.POST("/login/forgot", controller.Forgot)
			User.POST("/login/reset", controller.Reset)
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(os.Getenv("JWTSigningKey"))))
				Auth.Use(auth.AuthMiddleware)
				Auth.GET("/sessions", controller.Sessions)
				Auth.DELETE("/logout", controller.Logout)
			}
		}
	}

	Route.GET("/swagger/*", echoSwagger.WrapHandler)

	return Route
}
