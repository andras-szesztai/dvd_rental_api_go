package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorized(w, r, fmt.Errorf("authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			app.unauthorized(w, r, fmt.Errorf("invalid authorization header"))
			return
		}

		token := parts[1]
		if token == "" {
			app.unauthorized(w, r, fmt.Errorf("token is required"))
			return
		}

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["sub"].(float64)), 10, 64)
		if err != nil {
			app.unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		user, err := app.store.Users.GetUserByID(r.Context(), userId)
		if err != nil {
			app.unauthorized(w, r, fmt.Errorf("invalid token"))
			return
		}

		role, err := app.store.Roles.GetRoleByID(r.Context(), int64(user.Role.ID))
		if err != nil {
			app.unauthorized(w, r, fmt.Errorf("invalid token"))
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
			app.unauthorized(w, r, fmt.Errorf("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
