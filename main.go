package main

import (
	"os"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/config/database"

	"github.com/Gommunity/GoWithWith/config/mail"
	_ "github.com/Gommunity/GoWithWith/docs"
	"github.com/Gommunity/GoWithWith/routes"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/Gommunity/GoWithWith/services/validation"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zebresel-com/mongodm"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {

	utility.LoadEnvironmentVariables(".env")

	db := &database.Composer{
		Locals:   "resource/locals/locals.json",
		Addrs:    []string{os.Getenv("DBAddrs")},
		Database: os.Getenv("DBName"),
		Username: os.Getenv("DBUsername"),
		Password: os.Getenv("DBPassword"),
		Source:   os.Getenv("DBSource"),
	}
	db.Shoot(map[string]mongodm.IDocumentBase{
		"authAttempts": &model.AuthAttempt{},
		"sessions":     &model.Session{},
		"users":        &model.User{},
	})

	mail.Composer()
}

// @title GoWithWith
// @version 1.0
// @description A user system API starter.

// @contact.name Amir Irani
// @contact.email freshmanlimited@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3500
// @BasePath /endpoint

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
