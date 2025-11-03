package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeoService - мок для GeoService
type MockGeoService struct {
	mock.Mock
}

func (m *MockGeoService) AddressSearch(query string) ([]*Address, error) {
	args := m.Called(query)
	return args.Get(0).([]*Address), args.Error(1)
}

func (m *MockGeoService) GeoCode(lat, lng string) ([]*Address, error) {
	args := m.Called(lat, lng)
	return args.Get(0).([]*Address), args.Error(1)
}

func TestHandlerAddressSearch(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockQuery      string
		mockResponse   []*Address
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful address search",
			requestBody: RequestAddressSearch{
				Query: "Москва Ленина",
			},
			mockQuery: "Москва Ленина",
			mockResponse: []*Address{
				{
					City:   "Москва",
					Street: "Ленина",
					House:  "10",
					Lat:    "55.7558",
					Lon:    "37.6173",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{
					{
						City:   "Москва",
						Street: "Ленина",
						House:  "10",
						Lat:    "55.7558",
						Lon:    "37.6173",
					},
				},
			},
		},
		{
			name: "multiple addresses search",
			requestBody: RequestAddressSearch{
				Query: "Санкт-Петербург Невский",
			},
			mockQuery: "Санкт-Петербург Невский",
			mockResponse: []*Address{
				{
					City:   "Санкт-Петербург",
					Street: "Невский проспект",
					House:  "1",
					Lat:    "59.9390",
					Lon:    "30.3158",
				},
				{
					City:   "Санкт-Петербург",
					Street: "Невский проспект",
					House:  "25",
					Lat:    "59.9350",
					Lon:    "30.3258",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{
					{
						City:   "Санкт-Петербург",
						Street: "Невский проспект",
						House:  "1",
						Lat:    "59.9390",
						Lon:    "30.3158",
					},
					{
						City:   "Санкт-Петербург",
						Street: "Невский проспект",
						House:  "25",
						Lat:    "59.9350",
						Lon:    "30.3258",
					},
				},
			},
		},
		{
			name: "empty query",
			requestBody: RequestAddressSearch{
				Query: "",
			},
			mockQuery:      "",
			mockResponse:   []*Address{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{},
			},
		},
		{
			name: "service error",
			requestBody: RequestAddressSearch{
				Query: "invalid query",
			},
			mockQuery:      "invalid query",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   assert.AnError.Error(),
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"query": "test"`, // невалидный JSON
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty JSON object",
			requestBody:    map[string]interface{}{},
			mockQuery:      "",
			mockResponse:   []*Address{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок
			mockGeoService := new(MockGeoService)

			// Настраиваем мок только если ожидаем вызов
			if tt.mockQuery != "" || tt.mockResponse != nil {
				mockGeoService.On("AddressSearch", tt.mockQuery).Return(tt.mockResponse, tt.mockError)
			}

			// Создаем хендлер
			handler := HandlerAddressSearch(mockGeoService)

			// Подготавливаем тело запроса
			var requestBodyBytes []byte
			switch body := tt.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(body)
			default:
				requestBodyBytes, _ = json.Marshal(body)
			}

			// Создаем запрос
			req := httptest.NewRequest("POST", "/address/search", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызываем хендлер
			handler.ServeHTTP(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем тело ответа
			if tt.expectedStatus == http.StatusOK {
				var response ResponseAddress
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else if tt.expectedStatus >= http.StatusBadRequest {
				if expectedError, ok := tt.expectedBody.(string); ok {
					assert.Contains(t, rr.Body.String(), expectedError)
				}
			}

			// Проверяем вызовы мока
			mockGeoService.AssertExpectations(t)
		})
	}
}

func TestHandlerAddressGeocode(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockLat        string
		mockLng        string
		mockResponse   []*Address
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful geocode",
			requestBody: RequestAddressGeocode{
				Lat: "55.7558",
				Lng: "37.6173",
			},
			mockLat: "55.7558",
			mockLng: "37.6173",
			mockResponse: []*Address{
				{
					City:   "Москва",
					Street: "Красная площадь",
					House:  "1",
					Lat:    "55.7539",
					Lon:    "37.6208",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{
					{
						City:   "Москва",
						Street: "Красная площадь",
						House:  "1",
						Lat:    "55.7539",
						Lon:    "37.6208",
					},
				},
			},
		},
		{
			name: "geocode with multiple addresses",
			requestBody: RequestAddressGeocode{
				Lat: "59.9390",
				Lng: "30.3158",
			},
			mockLat: "59.9390",
			mockLng: "30.3158",
			mockResponse: []*Address{
				{
					City:   "Санкт-Петербург",
					Street: "Дворцовая площадь",
					House:  "2",
					Lat:    "59.9390",
					Lon:    "30.3158",
				},
				{
					City:   "Санкт-Петербург",
					Street: "Дворцовая набережная",
					House:  "38",
					Lat:    "59.9410",
					Lon:    "30.3178",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{
					{
						City:   "Санкт-Петербург",
						Street: "Дворцовая площадь",
						House:  "2",
						Lat:    "59.9390",
						Lon:    "30.3158",
					},
					{
						City:   "Санкт-Петербург",
						Street: "Дворцовая набережная",
						House:  "38",
						Lat:    "59.9410",
						Lon:    "30.3178",
					},
				},
			},
		},
		{
			name: "empty coordinates",
			requestBody: RequestAddressGeocode{
				Lat: "",
				Lng: "",
			},
			mockLat:        "",
			mockLng:        "",
			mockResponse:   []*Address{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{},
			},
		},
		{
			name: "service error",
			requestBody: RequestAddressGeocode{
				Lat: "invalid",
				Lng: "invalid",
			},
			mockLat:        "invalid",
			mockLng:        "invalid",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   assert.AnError.Error(),
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"lat": "55.7558"`, // невалидный JSON
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing lng field",
			requestBody: map[string]interface{}{
				"lat": "55.7558",
			},
			mockLat:        "55.7558",
			mockLng:        "",
			mockResponse:   []*Address{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: ResponseAddress{
				Addresses: []*Address{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок
			mockGeoService := new(MockGeoService)

			// Настраиваем мок только если ожидаем вызов
			if tt.mockLat != "" || tt.mockLng != "" || tt.mockResponse != nil {
				mockGeoService.On("GeoCode", tt.mockLat, tt.mockLng).Return(tt.mockResponse, tt.mockError)
			}

			// Создаем хендлер
			handler := HandlerAddressGeocode(mockGeoService)

			// Подготавливаем тело запроса
			var requestBodyBytes []byte
			switch body := tt.requestBody.(type) {
			case string:
				requestBodyBytes = []byte(body)
			default:
				requestBodyBytes, _ = json.Marshal(body)
			}

			// Создаем запрос
			req := httptest.NewRequest("POST", "/address/geocode", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Создаем ResponseRecorder
			rr := httptest.NewRecorder()

			// Вызываем хендлер
			handler.ServeHTTP(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем тело ответа
			if tt.expectedStatus == http.StatusOK {
				var response ResponseAddress
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else if tt.expectedStatus >= http.StatusBadRequest {
				if expectedError, ok := tt.expectedBody.(string); ok {
					assert.Contains(t, rr.Body.String(), expectedError)
				}
			}

			// Проверяем вызовы мока
			mockGeoService.AssertExpectations(t)
		})
	}
}

func TestRoutes(t *testing.T) {
	// Создаем мок сервиса
	mockGeoService := new(MockGeoService)
	mockGeoService.On("AddressSearch", "test").Return([]*Address{
		{
			City:   "Москва",
			Street: "Тестовая",
			House:  "1",
			Lat:    "55.7558",
			Lon:    "37.6173",
		},
	}, nil)
	mockGeoService.On("GeoCode", "55.7558", "37.6173").Return([]*Address{
		{
			City:   "Москва",
			Street: "Кремль",
			House:  "1",
			Lat:    "55.7539",
			Lon:    "37.6208",
		},
	}, nil)

	// Создаем роутер как в main()
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Post("/address/search", HandlerAddressSearch(mockGeoService))
		r.Post("/address/geocode", HandlerAddressGeocode(mockGeoService))
	})

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:   "address search route",
			method: "POST",
			path:   "/api/address/search",
			body: RequestAddressSearch{
				Query: "test",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "geocode route",
			method: "POST",
			path:   "/api/address/geocode",
			body: RequestAddressGeocode{
				Lat: "55.7558",
				Lng: "37.6173",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found route",
			method:         "GET",
			path:           "/api/unknown",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "wrong method for search",
			method:         "GET",
			path:           "/api/address/search",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBodyBytes []byte
			if tt.body != nil {
				requestBodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(requestBodyBytes))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

// Дополнительные тесты для проверки структуры ответа
func TestResponseStructure(t *testing.T) {
	mockGeoService := new(MockGeoService)
	mockGeoService.On("AddressSearch", "test").Return([]*Address{
		{
			City:   "Москва",
			Street: "Тестовая",
			House:  "1",
			Lat:    "55.7558",
			Lon:    "37.6173",
		},
	}, nil)

	handler := HandlerAddressSearch(mockGeoService)

	// Тестируем структуру ответа
	req := httptest.NewRequest("POST", "/address/search",
		bytes.NewBufferString(`{"query": "test"}`))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Проверяем Content-Type
	handler.ServeHTTP(rr, req)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Проверяем структуру JSON ответа
	var response ResponseAddress
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Addresses, 1)
	assert.Equal(t, "Москва", response.Addresses[0].City)
	assert.Equal(t, "Тестовая", response.Addresses[0].Street)
	assert.Equal(t, "1", response.Addresses[0].House)
	assert.Equal(t, "55.7558", response.Addresses[0].Lat)
	assert.Equal(t, "37.6173", response.Addresses[0].Lon)
}

// Тест на проверку корректности JSON маршалинга
func TestAddressJSONMarshaling(t *testing.T) {
	address := &Address{
		City:   "Москва",
		Street: "Тверская",
		House:  "15",
		Lat:    "55.7600",
		Lon:    "37.6100",
	}

	jsonBytes, err := json.Marshal(address)
	assert.NoError(t, err)

	var unmarshaledAddress Address
	err = json.Unmarshal(jsonBytes, &unmarshaledAddress)
	assert.NoError(t, err)
	assert.Equal(t, address, &unmarshaledAddress)
}
