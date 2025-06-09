package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()

	t.Run("should return bad request if payload is invalid", func(t *testing.T) {
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
}
