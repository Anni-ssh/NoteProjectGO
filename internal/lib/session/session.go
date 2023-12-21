package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// SetSession создаёт сессию через структуру http.Cookie.
func SetSession(sessionName string, tokenSize int, timeHour int, cookie *http.Cookie) error {
	const OP = "session.SetSessionOneDay"

	cookie.Name = sessionName

	token, err := generateRandomHash(tokenSize)
	if err != nil {
		return fmt.Errorf("session creation error: %w. Operation %s", err, OP)
	}
	cookie.Value = token

	expiration := time.Now().Add(time.Duration(timeHour) * time.Hour)
	cookie.Expires = expiration

	return nil
}

// generateRandomHash создает случайный хэш указанного размера.
func generateRandomHash(size int) (string, error) {
	// Создаем буфер для случайных байт
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Преобразуем случайные байты в строку в формате шестнадцатеричного представления
	randomHash := hex.EncodeToString(bytes)

	return randomHash, nil
}
