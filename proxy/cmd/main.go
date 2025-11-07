package main

import (
	"log"
	"net/http"

	_ "test/docs"

	"test/config"
	"test/internal/auth"
	"test/internal/geo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title HugoProxy Address API with JWT Auth
// @version 1.0
// @description API для поиска адресов и геокодирования с JWT аутентификацией
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/mihhha985/hogoproxy
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	var config = config.LoadConfig()
	var user = &auth.User{}
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(config.JwtSecret), nil)
	geoService := geo.NewGeoService(config.DaDataAPIKey, config.DaDataSecretKey)
	geoController := geo.NewGeoController(geoService, tokenAuth)
	authController := auth.NewAuthController(tokenAuth, user)

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authController.Register())
		r.Post("/login", authController.Login())

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			// Protected routes
			r.Post("/address/search", geoController.HandlerAddressSearch())
			r.Post("/address/geocode", geoController.HandlerAddressGeocode())
		})
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
