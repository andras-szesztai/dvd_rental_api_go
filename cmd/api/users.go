package main

import (
	"fmt"
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
)

type RegisterAdminUserPayload struct {
	Username  string `json:"username" validate:"required,min=3,max=20"`
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	StoreID   int    `json:"store_id" validate:"required,min=1"`
}

// Only authenticated Admins should be able to create
func (app *application) registerAdminUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterAdminUserPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = Validator.Struct(payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	rentalPlace, err := app.store.RentalPlace.GetRentalPlaceByID(r.Context(), payload.StoreID)
	if err != nil {
		app.badRequest(w, r, fmt.Errorf("failed to get rental place: %w", err))
		return
	}

	user := store.AdminUser{
		AddressID: rentalPlace.ID,
		Username:  payload.Username,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		StoreID:   payload.StoreID,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.User.CreateAdminUser(r.Context(), &user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, nil)
}
