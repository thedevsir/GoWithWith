package routes

import (
	"os"

	"github.com/Gommunity/GoWithWith/auth"
	"github.com/Gommunity/GoWithWith/routes/handlers"
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
			User.POST("/signup", handlers.Signup)
			User.POST("/signup/resend-email", handlers.ResendEmail)
			User.POST("/login", handlers.Login)
			User.POST("/login/forgot", handlers.Forgot)
			User.POST("/login/reset", handlers.Reset)
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(os.Getenv("JWTSigningKey"))))
				Auth.Use(auth.AuthMiddleware)
				Auth.GET("/sessions", handlers.Sessions)
				Auth.DELETE("/logout", handlers.Logout)
			}
		}
	}

	Route.GET("/swagger/*", echoSwagger.WrapHandler)

	return Route
}
