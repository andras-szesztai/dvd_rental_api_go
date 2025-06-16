package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"github.com/go-chi/chi/v5"
)

type rentalResponse struct {
	Data store.Rental `json:"data"`
}

// GetRentalByID godoc
//
//	@Summary		Get rental by ID
//	@Description	Get a rental by ID
//	@Tags			4. Rentals
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Rental ID"
//	@Success		200	{object}	rentalResponse
//	@Failure		400	{object}	utils.ErrorResponse
//	@Failure		401	{object}	utils.ErrorResponse
//	@Failure		404	{object}	utils.ErrorResponse
//	@Failure		500	{object}	utils.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/rentals/{id} [get]
func (app *application) getRentalByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := app.getUserContext(r)

	fmt.Println("user.ID", user.ID)

	if user == nil {
		app.errorHandler.Unauthorized(w, r, errors.New("unauthorized"))
		return
	}

	rentalID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	rental, err := app.store.Rentals.GetRental(r.Context(), rentalID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorHandler.NotFound(w, r)
			return
		}
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

	if user.Role.Name != "admin" {
		app.errorHandler.Unauthorized(w, r, errors.New("unauthorized"))
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, rentalResponse{Data: *rental}); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
	}
}
