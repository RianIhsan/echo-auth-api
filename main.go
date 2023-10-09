package main

import (
	"echo-auth-crud/config"
	"echo-auth-crud/router"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.Migration()

	app := echo.New()

	router.SetupRoute(app)

	app.Logger.Fatal(app.Start(":8080"))
}
