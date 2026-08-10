package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
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

// authConfig хранит пароль, прочитанный один раз при запуске сервера.
type authConfig struct {
	password string
}

func signInHandler(config authConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
			return
		}
		if config.password == "" {
			writeError(w, http.StatusNotFound, "аутентификация не настроена")
			return
		}

		var request signInRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "некорректный JSON")
			return
		}
		if subtle.ConstantTimeCompare([]byte(request.Password), []byte(config.password)) != 1 {
			writeError(w, http.StatusUnauthorized, "неверный пароль")
			return
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}).SignedString([]byte(config.password))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "не удалось создать токен")
			return
		}
		writeJSON(w, signInResponse{Token: token})
	}
}

func auth(config authConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.password == "" {
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
			return []byte(config.password), nil
		})
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
