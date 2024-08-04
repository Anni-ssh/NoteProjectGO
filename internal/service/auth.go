package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go
const (
	salt      = "fkdie293dkenr"
	tokensTTl = time.Hour * 24
)

type AuthService struct {
	storage storage.Authorization
}

func NewAuthService(s storage.Authorization) *AuthService {
	return &AuthService{storage: s}
}

func (s *AuthService) CreateUser(user entities.User) (int, error) {
	user.Password = genPasswordHash(user.Password)
	return s.storage.CreateUser(user)
}

func (s *AuthService) CheckUser(username, password string) (entities.User, error) {
	return s.storage.CheckUser(username, genPasswordHash(password))
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

func (s *AuthService) GenToken(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokensTTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.Id,
	})
	return token.SignedString([]byte(salt))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(salt), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserID, nil
}
