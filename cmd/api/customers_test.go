package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()

	token, err := app.authenticator.GenerateToken(jwt.MapClaims{})
	assert.NoError(t, err)

	t.Run("it should return bad request if payload is invalid", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID:   1,
					Name: "admin",
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{
				ID:   1,
				Name: "admin",
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/customers", bytes.NewBufferString(`{"store_id": 1, "first_name": "John", "last_name": "Doe"}`))
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'createCustomerPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag")
	})

	t.Run("customer cannot create customers", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID:   1,
					Name: "customer",
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{
				ID:   1,
				Name: "customer",
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/customers", bytes.NewBufferString(`{"store_id": 1, "first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}`))
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("only staff should be able to create customers (happy path)", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID:   2,
					Name: "admin",
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{
				ID:   2,
				Name: "admin",
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/customers", bytes.NewBufferString(`{"store_id": 1, "first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}`))
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.NotContains(t, recorder.Body.String(), "data")
	})

	t.Run("it returns internal server error if create failed", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID:   2,
					Name: "admin",
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{
				ID:   2,
				Name: "admin",
			}, nil
		}
		app.store.Customers.(*store.MockCustomerStore).CreateCustomerFunc = func(ctx context.Context, customer *store.Customer) error {
			return errors.New("create failed")
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/customers", bytes.NewBufferString(`{"store_id": 1, "first_name": "John", "last_name": "Doe", "email": "john.doe@example.com"}`))
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "the server encountered a problem and could not process your request")
	})
}
