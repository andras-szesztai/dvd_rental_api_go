package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
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
	currentUser := app.getUserContext(r)

	if currentUser.Role != "admin" {
		app.unauthorized(w, r, fmt.Errorf("unauthorized"))
		return
	}

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

	storeID, err := strconv.ParseInt(strconv.Itoa(payload.StoreID), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	rentalPlace, err := app.store.RentalPlaces.GetRentalPlaceByID(r.Context(), storeID)
	if err != nil {
		app.badRequest(w, r, fmt.Errorf("failed to get rental place: %w", err))
		return
	}

	newUser := store.User{
		AddressID: rentalPlace.ID,
		Username:  payload.Username,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		StoreID:   payload.StoreID,
	}

	if err := newUser.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Users.CreateAdminUser(r.Context(), &newUser)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, nil)
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (app *application) signInAdminUser(w http.ResponseWriter, r *http.Request) {
	var payload SignInPayload

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

	user, err := app.store.Users.GetAdminUserByEmail(r.Context(), payload.Email)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorized(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.aud,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

const userContextKey = contextKey("user")

func (app *application) getUserContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userContextKey).(*store.User)
	return user
}
