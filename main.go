package main

import (
	"os"

	db "./database"
	_ "./docs"
	mail "./gomail"
	models "./models"
	routes "./routes"
	handlers "./routes/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zebresel-com/mongodm"
)

func init() {

	// Load Environments
	err := godotenv.Load()
	if err != nil {
		panic(":main:init: ErrorLoading.EnvFile")
	}

	// Init Config in Handlers
	handlers.InitConfig()

	// Init Mail Driver
	mail.Initial()

	// Database Models
	Models := make(map[string]mongodm.IDocumentBase)
	Models["authAttempts"] = &models.AuthAttempt{}
	Models["sessions"] = &models.Session{}
	Models["users"] = &models.User{}

	// Setting up Database with Models
	db.Initial(Models, false)
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

	// Logger Middleware
	Run.Use(middleware.Logger())

	// Logger Middleware
	Run.Use(middleware.Recover())

	// Body Dump Middleware
	// Run.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))

	// Body Limit Middleware
	Run.Use(middleware.BodyLimit("10K"))

	// CORS Middleware
	Run.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Gzip Middleware
	// has conflict with swagger
	// Run.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 5,
	// }))

	// Recover Middleware
	// Run.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	// 	StackSize: 1 << 10, // 1 KB
	// }))

	// HTTPS NonWWW Redirect
	// Run.Pre(middleware.HTTPSNonWWWRedirect())

	// Start Server
	Run.Logger.Fatal(Run.Start(os.Getenv("PORT")))
}
