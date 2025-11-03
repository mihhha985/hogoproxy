package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "test/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"golang.org/x/crypto/bcrypt"
)

const DADATA_API_KEY = "ced67ee66aaf9f6df09e8e17e7ce3ffb56a05f8c"
const DADATA_SECRET_KEY = "d2ecbadfc616acaa12cbd48270e5fe685b8eb7fc"

// GeoServiceInterface defines the interface for geo services
type GeoServiceInterface interface {
	AddressSearch(query string) ([]*Address, error)
	GeoCode(lat, lng string) ([]*Address, error)
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server
// @host localhost:8080
// @BasePath /api
func main() {
	var user = &User{}
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	geoService := NewGeoService(DADATA_API_KEY, DADATA_SECRET_KEY)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", user.Register(tokenAuth))
		r.Post("/login", user.Login(tokenAuth))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			//r.Use(jwtauth.Authenticator(tokenAuth))

			// Protected routes
			r.Post("/address/search", HandlerAddressSearch(geoService))
			r.Post("/address/geocode", HandlerAddressGeocode(geoService))
		})
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func (user *User) Register(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var data User
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"email": data.Username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedBytes)
		user.Username = data.Username
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}

func (user *User) Login(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var data User
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"email": data.Username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if data.Username != user.Username {
			http.Error(w, "invalid username or password", http.StatusOK)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
			http.Error(w, "invalid username or password", http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}

// HandlerAddressSearch handles address search requests
// @Summary Search for addresses
// @Description Search for addresses using a text query
// @Tags address
// @Accept json
// @Produce json
// @Param request body RequestAddressSearch true "Search query"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /address/search [post]
func HandlerAddressSearch(geoService GeoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressSearch
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		log.Println("Authenticated user:", claims["email"])
		addresses, err := geoService.AddressSearch(req.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := ResponseAddress{Addresses: addresses}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// HandlerAddressGeocode handles geocoding requests
// @Summary Geocode coordinates to address
// @Description Convert latitude and longitude coordinates to address
// @Tags address
// @Accept json
// @Produce json
// @Param request body RequestAddressGeocode true "Coordinates"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /address/geocode [post]
func HandlerAddressGeocode(geoService GeoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressGeocode
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		log.Println("Authenticated user:", claims["email"])
		addresses, err := geoService.GeoCode(req.Lat, req.Lng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := ResponseAddress{Addresses: addresses}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
