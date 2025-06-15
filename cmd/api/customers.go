package main

import (
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type createCustomerPayload struct {
	StoreID   int64  `json:"store_id" validate:"required,min=1" example:"1"`
	FirstName string `json:"first_name" validate:"required,min=3,max=20" example:"John"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20" example:"Doe"`
	Email     string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

// CreateCustomer godoc
//
//	@Summary		Create customer
//	@Description	Create a new customer for store by admin user
//	@Tags			3. Customers
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createCustomerPayload	true	"Create customer request"
//	@Success		201		{object}	nil
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/customers [post]
func (app *application) createCustomer(w http.ResponseWriter, r *http.Request) {
	var payload createCustomerPayload

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

	customer := &store.Customer{
		StoreID:   payload.StoreID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	err = app.store.Customers.CreateCustomer(r.Context(), customer)
	if err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusCreated, nil); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}
}
