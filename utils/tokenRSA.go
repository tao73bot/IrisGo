package utils

// import (
// 	"log"
// 	"time"

// 	"github.com/kataras/iris/v12/middleware/jwt"
// )

// const (
// 	accessTokenMaxAge  = 10 * time.Minute
// 	refreshTokenMaxAge = 2 * time.Hour
// )

// var (
// 	privateKey, publicKey = jwt.MustLoadRSA("rsa_private_key.pem", "rsa_public_key.pem")

// 	signer    = jwt.NewSigner(jwt.RS256, privateKey, accessTokenMaxAge)
// 	refSigner = jwt.NewSigner(jwt.RS256, privateKey, refreshTokenMaxAge)
// 	verifer   = jwt.NewVerifier(jwt.RS256, publicKey)
// )

// type SignedDetails struct {
// 	Email  string
// 	Name   string
// 	Role   string
// 	UserId string
// }

// func GenerateAllToken(email string, name string, role string, userId string) (signedToken string, signedRefresh string, err error) {
// 	claims := &SignedDetails{
// 		Email:  email,
// 		Name:   name,
// 		Role:   role,
// 		UserId: userId,
// 	}
// 	refreshClaims := &SignedDetails{
// 		Email:  email,
// 		Name:   name,
// 		Role:   role,
// 		UserId: userId,
// 	}
// 	token, err := signer.Sign(claims)
// 	if err != nil {
// 		log.Fatal("Error in generating access token")
// 	}
// 	refreshToken, err := refSigner.Sign(refreshClaims)
// 	if err != nil {
// 		log.Fatal("Error in generating refresh token")
// 	}
// 	return string(token), string(refreshToken), nil
// }

// func ValidateTokens(tokenString string) (claims *SignedDetails, err error) {
// 	claims = &SignedDetails{}
// 	verifer.WithDefaultBlocklist()

// 	// Verify the token
// 	err = verifer.Verify([]byte(tokenString), claims)
// 	if err != nil {
// 		return
// 	}

// 	// Check token expiration
// 	currentTime := time.Now()
// 	if currentTime.After(currentTime.Add(accessTokenMaxAge)) {
// 		return nil, jwt.ErrTokenExpired
// 	}

// 	return claims, nil
// }