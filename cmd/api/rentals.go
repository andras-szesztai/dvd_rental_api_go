package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
