package main

import (
	"os"

	"github.com/Gommunity/GoWithWith/app/controller"
	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/config/mail"
	"github.com/Gommunity/GoWithWith/helpers/validation"
	"github.com/Gommunity/GoWithWith/models"
	"github.com/Gommunity/GoWithWith/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zebresel-com/mongodm"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		panic(":main:init: ErrorLoading.EnvFile")
	}

	controller.InitConfig()
	mail.Initial()

	Models := map[string]mongodm.IDocumentBase{
		"authAttempts": &models.AuthAttempt{},
		"sessions":     &models.Session{},
		"users":        &models.User{},
	}
	database.Initial(Models, false)
}

// @title GoWithWith
// @version 1.0
// @description A user system API starter.

// @contact.name Amir Irani
// @contact.email freshmanlimited@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3500
// @BasePath /endpoint/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	Run := routes.Initial()
	Run.Use(middleware.Logger())
	Run.Use(middleware.Recover())
	Run.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))
	Run.Use(middleware.BodyLimit("10K"))
	Run.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	Run.Validator = &validation.DataValidator{ValidatorData: validator.New()}

	// Gzip Middleware
	// has conflict with swagger
	// Run.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 5,
	// }))

	// Run.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	// 	StackSize: 1 << 10, // 1 KB
	// }))

	// Run.Pre(middleware.HTTPSNonWWWRedirect())

	Run.Logger.Fatal(Run.Start(os.Getenv("PORT")))
}
