package geo

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
