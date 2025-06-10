package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()

	t.Run("it should return bad request if payload is invalid", func(t *testing.T) {
		// no username
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Username' Error:Field validation for 'Username' failed on the 'required' tag")

		// request with invalid email
		req, err = http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Email' Error:Field validation for 'Email' failed on the 'email' tag")

		// no email
		req, err = http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag")

		// request with invalid password
		req, err = http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "pass"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Password' Error:Field validation for 'Password' failed on the 'min' tag")

		// no password
		req, err = http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Password' Error:Field validation for 'Password' failed on the 'required' tag")

		// request with invalid username
		req, err = http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "te", "password": "password"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'registerUserPayload.Username' Error:Field validation for 'Username' failed on the 'min' tag")
	})

	t.Run("it should return bad request if staff member already registered", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return &store.Staff{
				UserID: &[]int{1}[0],
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "staff member already registered")
	})

	t.Run("it should return bad request if staff member lookup returns error that is not sql.ErrNoRows", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "some error")
	})

	// should go ahead if sql no rows for staff
	// should assign admin role if staff member is found and not registered
	// it should return bad request if staff member lookup returns error, even if sql.ErrNoRows
	// it should return bad request if customer is already registered
	// should assign customer role if customer is found and not registered
	// bad request if role lookup returns error
	// bad request if user registration returns error
	// success if user registration returns no error, expected 201/

	//
}
