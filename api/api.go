package api

import (
	"Scrunchy/controllers"
	"Scrunchy/middleware"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{addr: addr}
}

func (s *ApiServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /products", controllers.GetProducts)
	router.HandleFunc("GET /products/{id}", controllers.GetProductIndex)

	router.HandleFunc("POST /users/signup", controllers.SignUp)
	router.HandleFunc("POST /users/login", controllers.Login)

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("GET /users/logout", controllers.Logout)
	authRouter.HandleFunc("GET /users/validate", controllers.Validate)

	authRouter.HandleFunc("GET /products/cart/", controllers.ListCart)
	authRouter.HandleFunc("POST /products/cart/add", controllers.AddToCart)
	authRouter.HandleFunc("POST /products/cart/remove", controllers.RemoveFromCart)
	authRouter.HandleFunc("POST /products/cart/quantity", controllers.QuantityCart)

	authRouter.HandleFunc("GET /admin/isadmin", middleware.RequireAdmin(http.HandlerFunc(controllers.IsAdmin)))
	authRouter.HandleFunc("POST /admin/products", middleware.RequireAdmin(http.HandlerFunc(controllers.PostProduct)))
	authRouter.HandleFunc("POST /admin/status", middleware.RequireAdmin(http.HandlerFunc(controllers.AdminStatus)))
	authRouter.HandleFunc("GET /admin/listusers", middleware.RequireAdmin(http.HandlerFunc(controllers.ListUsers)))

	router.Handle("/", middleware.RequireAuth(authRouter))

	stack := middleware.MiddlewareChain(middleware.Logger, middleware.RecoveryMiddleware)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Specify your React frontend origin
		AllowCredentials: true,                              // Allow cookies and credentials
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(stack(router))

	server := http.Server{
		Addr:    s.addr,
		Handler: corsHandler,
	}
	fmt.Println("Server has started", s.addr)
	return server.ListenAndServe()
}
