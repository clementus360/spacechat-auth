package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(userId string, jwtSecret string) (string, error) {
	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new JWT token and sign it with the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string, jwtSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Failed to extract the claims")
	}

	userId, ok := (*claims)["sub"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to extract userId from claims")
	}

	return userId, nil
}
