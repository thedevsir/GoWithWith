package routes

import (
	"os"

	auth "../auth"
	handler "./handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
)

// Initial ...
func Initial() *echo.Echo {

	Route := echo.New()

	endpointV1 := Route.Group("/endpoint/v1")
	{
		User := endpointV1.Group("/user")
		{
			User.POST("/signup", handler.Signup)
			User.POST("/login", handler.Login)
			User.POST("/forgot", handler.Forgot)
			User.POST("/reset", handler.Reset)
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(os.Getenv("JWTSigningKey"))))
				Auth.Use(auth.AuthMiddleware)
				Auth.GET("/sessions", handler.Sessions)
				Auth.DELETE("/logout", handler.Logout)
			}
		}
	}

	Route.GET("/swagger/*", echoSwagger.WrapHandler)

	return Route
}
