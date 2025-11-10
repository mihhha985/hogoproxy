package controller

import "test/internal/service"

// GeoServiceInterface defines the interface for geo services
type GeoServiceInterface interface {
	AddressSearch(query string) ([]*service.Address, error)
	GeoCode(lat, lng string) ([]*service.Address, error)
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []*service.Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
