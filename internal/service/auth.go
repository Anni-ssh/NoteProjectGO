package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"NoteProject/internal/storage"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

type AuthService struct {
	storage storage.Authorization
}

func NewAuthService(s storage.Authorization) *AuthService {
	return &AuthService{storage: s}
}

func (s *AuthService) CreateUser(ctx context.Context, user entities.User) (int, error) {
	user.Password = genPasswordHash(user.Password)
	return s.storage.CreateUser(ctx, user)
}

func (s *AuthService) CheckUser(ctx context.Context, username, password string) (entities.User, error) {
	return s.storage.CheckUser(ctx, username, genPasswordHash(password))
}

// Время жизни сессии
const (
	tokensTTL = time.Hour * 24
)

// Salt для хеширования
var salt = os.Getenv("HASH_SALT")

// Генерация хеша пароля
func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type tokenClaims struct {
	UserID int
	jwt.StandardClaims
}

// Генерация токена
func (s *AuthService) GenToken(user entities.User) (string, error) {
	const op = "service.GenToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokensTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.Id,
	})

	result, err := token.SignedString([]byte(salt))
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

// Парсинг токена
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	const op = "ParseToken"

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(salt), nil
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, errs.ErrInvalidTokenType)
	}

	return claims.UserID, nil
}
