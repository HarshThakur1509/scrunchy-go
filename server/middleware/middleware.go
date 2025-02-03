package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"github.com/markbates/goth/gothic"
)

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
		// Retrieve user ID from the session
		userID, err := gothic.GetFromSession("user_id", r)
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Refresh the session's lifetime
		session, _ := gothic.Store.Get(r, gothic.SessionName)
		session.Options.MaxAge += 86400 // Extend by 1 day (adjust as needed)
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Failed to refresh session", http.StatusInternalServerError)
			return
		}

		// Attach user ID to the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAdmin(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Retrieve user ID from the session
		userID, err := gothic.GetFromSession("user_id", r)
		if err != nil || userID == "" {
			// Return an empty JSON response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]interface{}{"exists": false})
			return
		}

		var user models.User
		initializers.DB.Select("admin").First(&user, userID)

		// Check for admin status
		if !user.Admin {
			http.Error(w, "Forbidden - User is not an admin", http.StatusForbidden)
			return
		}

		// Continue processing the request if user is admin
		next.ServeHTTP(w, r)
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
