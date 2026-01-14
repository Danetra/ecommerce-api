package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID int, username string) (string, time.Time, error) {
	// expired 1 jam
	expiredAt := time.Now().Add(1 * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiredAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiredAt, nil
}
