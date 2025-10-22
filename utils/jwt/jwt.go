package jwt

import (
	"context"
	"easy-attend-service/models"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}

func VerifyToken(raw string) (map[string]any, error) {
	godotenv.Load()
	token, err := jwt.Parse(raw, func(token *jwt.Token) (
		interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token singing method")
		}
		secret := []byte(os.Getenv("JWT_SECRET"))
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func GenerateTokenTeacher(ctx context.Context, teacher *models.Teacher) (string, error) {
	godotenv.Load()

	// Get JWT expiry hours from environment, default to 24 hours
	expireHours := 24
	if hours := os.Getenv("JWT_EXPIRE_HOURS"); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			expireHours = h
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": jwt.MapClaims{
			"id":         teacher.ID,
			"email":      teacher.Email,
			"first_name": teacher.FirstName,
			"last_name":  teacher.LastName,
			"phone":      teacher.Phone,
		},
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("[error]: %v", err)
		return "", err
	}
	return tokenString, nil
}

func GenerateToken(claims CustomClaims) (string, time.Time, error) {
	godotenv.Load()

	// Get JWT expiry hours from environment, default to 24 hours
	expireHours := 24
	if hours := os.Getenv("JWT_EXPIRE_HOURS"); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			expireHours = h
		}
	}

	expiresAt := time.Now().Add(time.Duration(expireHours) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   claims.UserID,
		"email":     claims.Email,
		"user_type": claims.UserType,
		"nbf":       time.Now().Unix(),
		"exp":       expiresAt.Unix(),
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("[error]: %v", err)
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}
