package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/cache"
	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
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

	t.Run("it should assign admin role if staff member is found and not registered", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return &store.Staff{
				UserID: nil,
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByNameFunc = func(ctx context.Context, name string) (*store.Role, error) {
			return &store.Role{
				ID:   1,
				Name: "admin",
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("it should go ahead if sql no rows for staff but it should return bad request if customer lookup returns error, even if sql.ErrNoRows", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "failed to get customer")
	})

	t.Run("it should return bad request if customer is already registered", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return &store.Customer{
				UserID: &[]int{1}[0],
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "customer already registered")
	})

	t.Run("bad request if role lookup returns error", func(t *testing.T) {
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return &store.Customer{
				UserID: nil,
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByNameFunc = func(ctx context.Context, name string) (*store.Role, error) {
			return nil, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "some error")
	})

	t.Run("bad request if user registration returns error", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).RegisterUserFunc = func(ctx context.Context, user *store.User) error {
			return errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email": "test@test.com", "username": "test", "password": "password"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "some error")
	})
}

func TestSignInUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()
	plaintextPassword := "password"
	hashedPassword := utils.Password{Plaintext: &plaintextPassword}
	err := hashedPassword.Set(plaintextPassword)
	assert.NoError(t, err)

	t.Run("it should return bad request if payload is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "dasdas", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'signInPayload.Email' Error:Field validation for 'Email' failed on the 'email' tag")

		req, err = http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"password": "password"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'signInPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag")

		req, err = http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com"}`))
		assert.NoError(t, err)

		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Key: 'signInPayload.Password' Error:Field validation for 'Password' failed on the 'required' tag")
	})

	t.Run("bad request if user lookup returns error", func(t *testing.T) {
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com","password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "some error")
	})

	t.Run("no bad request if user lookup returns no rows and staff returns user id (happy staff path)", func(t *testing.T) {
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return &store.Staff{
				UserID: &[]int{1}[0],
			}, nil
		}
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID:       1,
				Password: hashedPassword,
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("bad request if staff lookup returns error (even no rows)", func(t *testing.T) {
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "sql: no rows in result set")
	})

	t.Run("bad request if user lookup returns error (even no rows)", func(t *testing.T) {
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return &store.Staff{
				UserID: &[]int{1}[0],
			}, nil
		}
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "sql: no rows in result set")
	})

	t.Run("bad request if password is incorrect", func(t *testing.T) {
		wrongPassword := "wrong"
		hashedPassword := utils.Password{Plaintext: &wrongPassword}
		err := hashedPassword.Set(wrongPassword)
		assert.NoError(t, err)
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return &store.Staff{
				UserID: &[]int{1}[0],
			}, nil
		}
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID:       1,
				Password: hashedPassword,
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("otherwise happy path for customer as well (happy customer path)", func(t *testing.T) {
		app.store.Customers.(*store.MockCustomerStore).GetCustomerByEmailFunc = func(ctx context.Context, email string) (*store.Customer, error) {
			return &store.Customer{
				UserID: &[]int{1}[0],
			}, nil
		}
		app.store.Staff.(*store.MockStaffStore).GetStaffByEmailFunc = func(ctx context.Context, email string) (*store.Staff, error) {
			return nil, sql.ErrNoRows
		}
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID:       1,
				Password: hashedPassword,
			}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/sign-in", bytes.NewBufferString(`{"email": "test@test.com", "password":"`+plaintextPassword+`"}`))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "data")
	})
}

func TestAuthTokenMiddleware(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mountRoutes()

	token, err := app.authenticator.GenerateToken(jwt.MapClaims{})
	assert.NoError(t, err)

	t.Run("it should return unauthorized if no authorization header is provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("it should return unauthorized if authorization header is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer%s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("it should return unauthorized if token is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer%s", "invalid"))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("unauthorized if user is not found", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "user not found")
	})

	t.Run("unauthorized if user has invalid role", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID:   12,
					Name: "invalid",
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{}, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "unauthorized")
	})

	t.Run("it reads user from cache if it exists", func(t *testing.T) {
		app.store.Users.(*store.MockUserStore).GetUserByIDFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{}, nil
		}
		app.cache.Users.(*cache.MockUserCache).GetFunc = func(ctx context.Context, id int64) (*store.User, error) {
			return &store.User{
				ID: 1,
				Role: &store.Role{
					ID: 1,
				},
			}, nil
		}
		app.store.Roles.(*store.MockRoleStore).GetRoleByIDFunc = func(ctx context.Context, id int64) (*store.Role, error) {
			return &store.Role{
				ID:   1,
				Name: "admin",
			}, nil
		}
		app.store.Rentals.(*store.MockRentalStore).GetRentalFunc = func(ctx context.Context, id int64) (*store.Rental, error) {
			return &store.Rental{
				ID: 1,
			}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/rentals/1", nil)
		assert.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
