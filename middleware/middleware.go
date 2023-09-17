package middleware

import (
	"echo-auth-crud/config"
	"echo-auth-crud/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or missing token"})
		}

		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Empty token"})
		}

		var user models.User
		if err := config.DB.Where("access_token = ?", token).First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		c.Set("user", &user)

		return next(c)
	}
}
