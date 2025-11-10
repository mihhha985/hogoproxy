package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"test/internal/responder"

	"github.com/go-chi/jwtauth/v5"
)

type GeoController struct {
	GeoService GeoServiceInterface
	TokenAuth  *jwtauth.JWTAuth
}

func NewGeoController(geoService GeoServiceInterface, tokenAuth *jwtauth.JWTAuth) *GeoController {
	return &GeoController{
		GeoService: geoService,
		TokenAuth:  tokenAuth,
	}
}

// HandlerAddressSearch handles address search requests
// @Summary Search for addresses
// @Description Search for addresses using a text query
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RequestAddressSearch true "Search query"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /address/search [post]
func (c *GeoController) HandlerAddressSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressSearch
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		log.Println("Authenticated user:", claims["email"])
		addresses, err := c.GeoService.AddressSearch(req.Query)
		if err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		resp := ResponseAddress{Addresses: addresses}
		responder.OutputJSON(w, resp)
	}
}

// HandlerAddressGeocode handles geocoding requests
// @Summary Geocode coordinates to address
// @Description Convert latitude and longitude coordinates to address
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RequestAddressGeocode true "Coordinates"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /address/geocode [post]
func (c *GeoController) HandlerAddressGeocode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressGeocode
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		log.Println("Authenticated user:", claims["email"])
		addresses, err := c.GeoService.GeoCode(req.Lat, req.Lng)
		if err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		resp := ResponseAddress{Addresses: addresses}
		responder.OutputJSON(w, resp)
	}
}
