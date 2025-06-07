package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getRentalByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rentalID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	rental, err := app.store.Rental.GetRental(r.Context(), rentalID)
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
