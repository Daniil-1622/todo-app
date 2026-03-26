package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt      = "jadnfIWFIAFIU123123ISJDNCIJSND"
	signedKey = "grjk#4%35FSFJlja#4353KSFjH"
	tokenTTL  = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{

			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signedKey))
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(signedKey)))
}

func (s *AuthService) ParseToken(accesstoken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		return []byte(signedKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	return claims.UserId, nil
}
