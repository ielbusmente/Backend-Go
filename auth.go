package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found, relying on system environment variables")
	}
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if jwtKey == nil {
		log.Fatal("JWT_SECRET_KEY is not set")
	}
}

func GenerateJWT(userID string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // expires in 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}
