package utils

import (
	"echo-auth-crud/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user models.User, secretKey string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   strconv.Itoa(int(user.Id)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))

}

func GetRefreshTokenSecret() string {
	return os.Getenv("REFRESH_TOKEN_SECRET")
}

func GetAccessTokenSecret() string {
	return os.Getenv("ACCESS_TOKEN_SECRET")
}

func VerifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
