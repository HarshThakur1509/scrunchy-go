package middleware

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SECRET = "l41^*&vjah4#%4565c4vty%#8b84"

type wrappedWrite struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWrite) writeHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWrite{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	}
}

func RecoveryMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack treace: %s"
				log.Printf(msg, err, string(debug.Stack()))

				er := http.StatusInternalServerError
				http.Error(w, "Internal Server Error", er)
			}
		}()

		next.ServeHTTP(w, r)
	}
}

func RequireAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Authorization cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and validate expiration
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - invalid claims", http.StatusUnauthorized)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			http.Error(w, "Unauthorized - token expired", http.StatusUnauthorized)
			return
		}

		// Retrieve user ID from claims and query database
		userID, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Unauthorized - invalid subject", http.StatusUnauthorized)
			return
		}
		var user models.User
		if err := initializers.DB.First(&user, userID).Error; err != nil || user.ID == 0 {
			http.Error(w, "Unauthorized - user not found", http.StatusUnauthorized)
			return
		}

		// Set user in the request context and pass it to the next handler
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAdmin(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Authorization cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and validate expiration
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - invalid claims", http.StatusUnauthorized)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			http.Error(w, "Unauthorized - token expired", http.StatusUnauthorized)
			return
		}

		// Retrieve user ID from claims and query database
		userID, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Unauthorized - invalid subject", http.StatusUnauthorized)
			return
		}
		var user models.User
		if err := initializers.DB.First(&user, userID).Error; err != nil || user.ID == 0 {
			http.Error(w, "Unauthorized - user not found", http.StatusUnauthorized)
			return
		}

		// Check if user is an admin
		if user.Admin {

			// Set user in the request context and pass it to the next handler
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized - User is not an Admin", http.StatusUnauthorized)
			return
		}

	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next.ServeHTTP
	}
}
