package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

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

	staff, err := app.store.Staff.GetStaffByEmail(r.Context(), payload.Email)
	if err != nil && err != sql.ErrNoRows {
		app.badRequest(w, r, fmt.Errorf("failed to get staff member: %w", err))
		return
	}

	if staff != nil && staff.UserID != nil {
		app.badRequest(w, r, fmt.Errorf("staff member already registered"))
		return
	}

	var roleName string
	if staff != nil && staff.ID > 0 {
		roleName = "admin"
	}

	if err == sql.ErrNoRows {
		customer, err := app.store.Customers.GetCustomerByEmail(r.Context(), payload.Email)
		if err != nil {
			app.badRequest(w, r, fmt.Errorf("failed to get customer: %w", err))
			return
		}

		if customer.UserID != nil {
			app.badRequest(w, r, fmt.Errorf("customer already registered"))
			return
		}

		roleName = "customer"
	}

	role, err := app.store.Roles.GetRoleByName(r.Context(), roleName)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &store.User{
		Email:    payload.Email,
		Username: payload.Username,
		Role:     role,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.store.Users.RegisterUser(r.Context(), user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type SignInPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (app *application) signInUser(w http.ResponseWriter, r *http.Request) {
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

	var userID int
	customer, err := app.store.Customers.GetCustomerByEmail(r.Context(), payload.Email)
	if err != nil && err != sql.ErrNoRows {
		app.badRequest(w, r, err)
		return
	}

	if customer == nil {
		staff, err := app.store.Staff.GetStaffByEmail(r.Context(), payload.Email)
		fmt.Println("staff error", staff)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		if staff.UserID != nil {
			userID = *staff.UserID
		}
	} else {
		userID = *customer.UserID
	}

	fmt.Println("customer error", userID)
	user, err := app.store.Users.GetUserByID(r.Context(), int64(userID))
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
