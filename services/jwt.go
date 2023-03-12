package services

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(userId string, jwtSecret string) (string, error){
	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new JWT token and sign it with the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err := token.SignedString([]byte(jwtSecret))
	if err!=nil {
		return "",nil
	}

	return tokenString,nil
}
