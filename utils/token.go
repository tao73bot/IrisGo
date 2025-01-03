package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var BlockList = make([]string, 200)
var CurrentIndex = 0

type SignedDetails struct {
	Email  string
	Name   string
	Role   string
	UserID string
	jwt.RegisteredClaims
}

func GenerateAllTokens(email string, name string, role string, userID string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:  email,
		Name:   name,
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	refreshClaims := &SignedDetails{
		Email:  email,
		Name:   name,
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal(err)
		return
	}
	return token, refreshToken, nil
}

func InvalidateToken(tokenString string) {
	for CurrentIndex < 200 {
		BlockList[CurrentIndex] = tokenString
		CurrentIndex = (CurrentIndex + 1) % 200
		break
	}
	// return nil
}

func ValidateToken(tokenString string) (claims *SignedDetails, err error) {
	claims = &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return
	}
	if _, ok := token.Claims.(*SignedDetails); !ok {
		return claims, jwt.ErrTokenInvalidClaims
	}
	if claims.ExpiresAt.Time.Unix() < time.Now().Local().Unix() {
		return claims, jwt.ErrTokenExpired
	}
	return claims, nil
}
