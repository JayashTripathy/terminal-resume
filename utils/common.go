package utils

import (
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"terminal-resume.jayash.space/models"
)

func IsValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseToken(tokenStr string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("my_secret"), nil
	} )

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}