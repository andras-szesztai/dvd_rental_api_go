package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetRentalByID godoc
//
//	@Summary		Get rental by ID
//	@Description	Get a rental by ID
//	@Tags			4. Rentals
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Rental ID"
//	@Success		200	{object}	rentalResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Security		ApiKeyAuth
//	@Router			/rentals/{id} [get]
func (app *application) getRentalByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := app.getUserContext(r)

	if user == nil || user.Role.Name != "admin" {
		app.unauthorized(w, r, errors.New("unauthorized"))
		return
	}

	rentalID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	rental, err := app.store.Rentals.GetRental(r.Context(), rentalID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFound(w, r)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, rentalResponse{Data: *rental})
}
