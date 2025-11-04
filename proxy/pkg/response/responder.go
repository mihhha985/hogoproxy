package response

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

//go:generate easytags $GOFILE
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

func OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Println("response writer error on write:", err)
	}
}

func ErrorBadRequest(w http.ResponseWriter, err error) {
	log.Println("http response bad request status code:", err)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	var resp = Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("response writer error on write:", err)
	}
}

func ErrorForbidden(w http.ResponseWriter, err error) {
	log.Println("http response forbidden:", err)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	var resp = Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("response writer error on write:", err)
	}
}

func ErrorUnauthorized(w http.ResponseWriter, err error) {
	log.Println("http response Unauthorized:", err)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	var resp = Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("response writer error on write:", err)
	}
}

func ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	log.Println("http response internal error:", err)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	var resp = Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("response writer error on write:", err)
	}
}
