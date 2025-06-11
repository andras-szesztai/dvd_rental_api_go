package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type registerUserPayload struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Username string `json:"username" validate:"required,min=3,max=20" example:"john.doe"`
	Password string `json:"password" validate:"required,min=8,max=72" example:"password123"`
}

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Register a new user
//	@Tags			2. Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		registerUserPayload	true	"Register user request"
//	@Success		201		{object}	nil
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/auth/register [post]
func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	err = Validator.Struct(payload)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	var roleName string
	staff, err := app.store.Staff.GetStaffByEmail(r.Context(), payload.Email)
	if err != nil && err != sql.ErrNoRows {
		app.errorHandler.BadRequest(w, r, fmt.Errorf("f	: %w", err))
		return
	}
	if staff != nil && staff.UserID != nil {
		app.errorHandler.BadRequest(w, r, fmt.Errorf("staff member already registered"))
		return
	}
	if staff != nil && staff.ID > 0 {
		roleName = "admin"
	}

	if err == sql.ErrNoRows {
		customer, err := app.store.Customers.GetCustomerByEmail(r.Context(), payload.Email)
		if err != nil {
			app.errorHandler.BadRequest(w, r, fmt.Errorf("failed to get customer: %w", err))
			return
		}
		if customer.UserID != nil {
			app.errorHandler.BadRequest(w, r, fmt.Errorf("customer already registered"))
			return
		}
		roleName = "customer"
	}

	role, err := app.store.Roles.GetRoleByName(r.Context(), roleName)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	user := &store.User{
		Email:    payload.Email,
		Username: payload.Username,
		Role:     role,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

	err = app.store.Users.RegisterUser(r.Context(), user)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusCreated, nil); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}
}

type signInPayload struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=8,max=72" example:"password123"`
}

type signInResponse struct {
	Data string `json:"data"`
}

// SignInUser godoc
//
//	@Summary		Sign in user
//	@Description	Sign in a user
//	@Tags			2. Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		signInPayload	true	"Sign in user request"
//	@Success		200		{object}	signInResponse	"JWT token"
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/auth/sign-in [post]
func (app *application) signInUser(w http.ResponseWriter, r *http.Request) {
	var payload signInPayload

	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	err = Validator.Struct(payload)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	var userID int
	customer, err := app.store.Customers.GetCustomerByEmail(r.Context(), payload.Email)
	if err != nil && err != sql.ErrNoRows {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	if customer == nil {
		staff, err := app.store.Staff.GetStaffByEmail(r.Context(), payload.Email)
		fmt.Println("staff error", staff)
		if err != nil {
			app.errorHandler.BadRequest(w, r, err)
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
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.errorHandler.Unauthorized(w, r, err)
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
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, signInResponse{Data: token}); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("invalid authorization header"))
			return
		}

		token := parts[1]
		if token == "" {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("token is required"))
			return
		}

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["sub"].(float64)), 10, 64)
		if err != nil {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		user, err := app.store.Users.GetUserByID(r.Context(), userId)
		if err != nil {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		role, err := app.store.Roles.GetRoleByID(r.Context(), int64(user.Role.ID))
		if err != nil {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		user.Role = role

		ctx := context.WithValue(r.Context(), contextKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) CheckAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextKey("user")).(*store.User)
		if user.Role.Name != "admin" {
			app.errorHandler.Unauthorized(w, r, fmt.Errorf("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

const userContextKey = contextKey("user")

func (app *application) getUserContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userContextKey).(*store.User)
	return user
}
