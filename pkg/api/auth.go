package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenLifetime = 8 * time.Hour

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func password() string {
	return os.Getenv("TODO_PASSWORD")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "метод не поддерживается"})
		return
	}
	if password() == "" {
		writeJSON(w, map[string]string{"error": "аутентификация не настроена"})
		return
	}

	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, map[string]string{"error": "некорректный JSON"})
		return
	}
	if subtle.ConstantTimeCompare([]byte(request.Password), []byte(password())) != 1 {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(password()))
	if err != nil {
		writeJSON(w, map[string]string{"error": "не удалось создать токен"})
		return
	}
	writeJSON(w, signInResponse{Token: token})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := password()
		if secret == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := ""
		if cookie, err := r.Cookie("token"); err == nil {
			tokenString = cookie.Value
		} else if authorization := r.Header.Get("Authorization"); strings.HasPrefix(authorization, "Bearer ") {
			tokenString = strings.TrimPrefix(authorization, "Bearer ")
		}
		if tokenString == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
