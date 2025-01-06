package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type SignedDetailsIris struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}

func GenerateTokenIris(signer *jwt.Signer, email, name, role, userID string) (signedToken string, err error) {
	claims := SignedDetailsIris{
		Email:  email,
		Name:   name,
		Role:   role,
		UserID: userID,
	}
	token, err := signer.Sign(claims)
	if err != nil {
		log.Fatal(err)
		return
	}
	return string(token), nil
}

func ValidateTokenIris(tokenString string) (claims *SignedDetailsIris, err error) {
	// Create a new verifier with the same secret key
	verifier := jwt.NewVerifier(jwt.HS256, []byte(os.Getenv("JWT_SECRET")))
	verifier.WithDefaultBlocklist()

	// Set token max age
	verifier.Extractors = []jwt.TokenExtractor{jwt.FromHeader} // extract from Authorization: Bearer <TOKEN>
	verifier.Validators = []jwt.TokenValidator{jwt.Expected{
		NotBefore: time.Now().Unix(),
	}} // validate token expiration

	// Parse and validate the token
	token, err := verifier.VerifyToken([]byte(tokenString))
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims
	claims = new(SignedDetailsIris)
	err = token.Claims(&claims)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}

	return claims, nil
}

type TokenBlockList struct {
	redisClient *redis.Client
}

func NewTokenBlocklist(redisAddr string) (tbl *TokenBlockList,err error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	tbl = &TokenBlockList{
		redisClient: client,
	}
	if tbl.redisClient == nil {
		return nil, fmt.Errorf("failed to connect to redis")
	}
	log.Println("Connected to redis")
	return tbl, nil
}

func (tb *TokenBlockList) InvalidateTokenIris(tokenString string, exp time.Duration) error {
	if tb.redisClient == nil {
		return fmt.Errorf("redis client is nil")
	}
	if tokenString == "" {
		return fmt.Errorf("token is empty")
	}
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := tb.redisClient.Set(ctx, "blacklist:"+tokenString, true, exp).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}
	log.Println("Token successfully invalidated")
	return nil
}

func (tb *TokenBlockList) IsTokenInvalidIris(tokenString string) bool {
	if tb.redisClient == nil {
		log.Println("redis client is nil")
		return true
	}
	if tokenString == "" {
		log.Println("token is empty")
		return true
	}
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := tb.redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
	if err != nil {
		log.Printf("failed to check if token is invalid: %v", err)
		return true
	}
	return exists > 0
}
