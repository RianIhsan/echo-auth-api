package router

import (
	"echo-auth-crud/controllers"
	"echo-auth-crud/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoute(app *echo.Echo) {
	app.POST("/register", controllers.Register)
	app.POST("/login", controllers.Login)
	app.POST("/refresh", controllers.RefreshToken)

	apps := app.Group("/auth")
	apps.Use(middleware.Middleware)
	apps.POST("/logout", controllers.Logout)
	apps.GET("/user", controllers.GetUser)

	apps.POST("/book", controllers.AddBook)
	apps.GET("/book", controllers.GetBooks)
	apps.GET("/book/:id", controllers.GetBook)
	apps.PATCH("/book/:id", controllers.UpdateBook)
	apps.DELETE("/book/:id", controllers.DeleteBook)
}
