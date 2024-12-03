package middleware

import (
	"blogAPI/internal/config"
	"log"
	"net/http"
	"strings"

	"context"

	"github.com/golang-jwt/jwt"
)

var User_UUID string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			log.Println("No authorization token provided")
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, http.ErrAbortHandler
			}
			return config.JWTConfig.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to extract claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		User_UUID, ok = claims["user_uuid"].(string)
		if !ok {
			log.Println("User UUID not found in token claims")
			http.Error(w, "User UUID not found in token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_uuid", User_UUID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
