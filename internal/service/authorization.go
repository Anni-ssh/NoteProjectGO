package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt      = "fkdie293dkenr"
	tokensTTl = time.Hour * 24
	signKey   = "sd329sdlkj239"
)

type AuthService struct {
	repository storage.Authorization
}

func NewAuthService(repository storage.Authorization) *AuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) CreateUser(user entities.User) (int, error) {
	user.Password = genPasswordHash(user.Password)
	return s.repository.CreateUser(user)
}

func (s *AuthService) CheckUser(username, password string) (*entities.User, error) {
	return s.repository.CheckUser(username, genPasswordHash(password))
}

func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type tokenClaims struct {
	UserID int
	jwt.StandardClaims
}

func (s *AuthService) GenAuthToken(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokensTTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.Id,
	})
	return token.SignedString([]byte(signKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not if type")
	}
	return claims.UserID, nil
}
