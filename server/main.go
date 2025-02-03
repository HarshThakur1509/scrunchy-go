package main

import (
	"log"
	"net/http"
	"os"

	"github.com/HarshThakur1509/scrunchy-go/api"
	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	initializers.ConnectDB()

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || callbackURL == "" {
		log.Println("Error: Google OAuth environment variables are not set in .env")
	}

	// Configure session store
	sessionSecret := os.Getenv("SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-session-secret"
	}

	store := sessions.NewCookieStore([]byte(sessionSecret))

	isProduction := os.Getenv("ENV") == "production"

	domain := os.Getenv("SESSION_COOKIE_DOMAIN")
	if domain == "" {
		domain = "localhost"
	}
	log.Println("domain: ", domain)

	store.Options = &sessions.Options{
		HttpOnly: true,
		Secure:   isProduction, // Enable secure cookies in production
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	// Configure Google provider
	goth.UseProviders(
		google.New(clientID, clientSecret, callbackURL, "email", "profile"),
	)
}

func main() {
	server := api.NewApiServer(":3000")
	server.Run()
}
