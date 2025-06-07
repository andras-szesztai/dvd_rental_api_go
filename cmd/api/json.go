package main

import (
	"encoding/json"
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSONResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSONResponse(w, status, &errorResponse{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSONResponse(w, status, data)
}

type rentalResponse struct {
	Data store.Rental `json:"data"`
}
