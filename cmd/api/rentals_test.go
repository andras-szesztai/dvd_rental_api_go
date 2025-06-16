package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetRental(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()

	token, err := app.authenticator.GenerateToken(jwt.MapClaims{})
	assert.NoError(t, err)

	t.Run("admin should be able to get a rental", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 1,
				},
			}, nil
		}
		app.store.Rentals.(*store.MockRentalStore).GetRentalFunc = func(ctx context.Context, id int64) (*store.Rental, error) {
			return &store.Rental{
				ID:         1,
				RentalDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "data")
		assert.Contains(t, recorder.Body.String(), "2021-01-01T00:00:00Z")
		assert.Contains(t, recorder.Body.String(), "1")
	})

	t.Run("non-admin user should not be able to get a rental", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 2,
				},
			}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("unauthorized if user is not authenticated (nil)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("bad request if id is not a number", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 1,
				},
			}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/notanumber", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "invalid syntax")
	})

	t.Run("not found if rental does not exist", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 1,
				},
			}, nil
		}
		app.store.Rentals.(*store.MockRentalStore).GetRentalFunc = func(ctx context.Context, id int64) (*store.Rental, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1000000", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "not found")
	})

	t.Run("internal server error if database error", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 1,
				},
			}, nil
		}
		app.store.Rentals.(*store.MockRentalStore).GetRentalFunc = func(ctx context.Context, id int64) (*store.Rental, error) {
			return nil, errors.New("database error")
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "the server encountered a problem and could not process your request")
	})
}
