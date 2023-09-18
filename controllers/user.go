package controllers

import (
	"echo-auth-crud/config"
	"echo-auth-crud/models"
	"echo-auth-crud/utils"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		return err
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	config.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.Email == user.Email {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"message": "Email already exists",
		})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := argon2id.CreateHash(user.Password, &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	})
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	config.DB.Create(&user)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c echo.Context) error {
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		return err
	}

	// Temukan pengguna berdasarkan email
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}

	// Periksa apakah password benar
	match, err := argon2id.ComparePasswordAndHash(user.Password, existingUser.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email"})
	}

	if !match {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid password"})
	}

	// Buat token
	accessToken, err := utils.CreateToken(user, utils.GetAccessTokenSecret())
	if err != nil {
		return err
	}

	refreshToken, err := utils.CreateToken(user, utils.GetRefreshTokenSecret())
	if err != nil {
		return err
	}

	existingUser.AccessToken = accessToken
	existingUser.RefreshToken = refreshToken
	config.DB.Save(&existingUser)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c echo.Context) error {
	type PayloadRefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}

	payload := PayloadRefreshToken{}

	if err := c.Bind(&payload); err != nil {
		return err
	}

	// Mendapatkan secret untuk refresh token
	refreshTokenSecret := utils.GetRefreshTokenSecret()

	// Memeriksa apakah refresh token valid
	claims, err := utils.VerifyToken(payload.RefreshToken, refreshTokenSecret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid refresh token"})
	}

	// Temukan pengguna berdasarkan ID dalam token
	var user models.User
	if err := config.DB.First(&user, claims["id"]).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Buat token baru
	accessToken, err := utils.CreateToken(user, utils.GetAccessTokenSecret())
	if err != nil {
		return err
	}

	refreshToken, err := utils.CreateToken(user, utils.GetRefreshTokenSecret())
	if err != nil {
		return err
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Token refreshed successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(c echo.Context) error {
	type PayloadRefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}

	payload := PayloadRefreshToken{}

	if err := c.Bind(&payload); err != nil {
		return err
	}

	// Mendapatkan secret untuk refresh token
	refreshTokenSecret := utils.GetRefreshTokenSecret()

	// Memeriksa apakah refresh token valid
	claims, err := utils.VerifyToken(payload.RefreshToken, refreshTokenSecret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid refresh token"})
	}

	// Temukan pengguna berdasarkan ID dalam token
	userID := claims["id"]
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Menghapus token
	user.AccessToken = ""
	user.RefreshToken = ""
	config.DB.Save(&user)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

func GetUser(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")

	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or missing token"})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Empty token"})
	}

	var user models.User
	if err := config.DB.Preload("Books").Where("access_token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
