package main

import (
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type moviesResponse struct {
	Data []*store.Movie `json:"data"`
}

// GetMovies godoc
//
//	@Summary		Get movies
//	@Description	Get all movies
//	@Tags			5. Movies
//	@Accept			json
//	@Produce		json
//	@Param			query	query		utils.MovieQuery	true	"Query parameters"
//	@Success		200		{object}	moviesResponse
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/movies [get]
func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {
	movieQuery := utils.MovieQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	movieQuery, err := movieQuery.Parse(r)
	if err != nil {
		app.errorHandler.BadRequest(w, r, err)
		return
	}

	movies, err := app.store.Movies.GetMovies(r.Context(), &movieQuery)
	if err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}

	for _, movie := range movies {
		inventories, err := app.store.Inventory.GetMovieInventory(r.Context(), movie.ID)
		if err != nil {
			app.errorHandler.InternalServerError(w, r, err)
			return
		}
		for _, inventory := range inventories {
			inventoryID := int64(inventory.ID)
			rentals, err := app.store.Rentals.GetMovieRentals(r.Context(), inventoryID)
			if err != nil {
				app.errorHandler.InternalServerError(w, r, err)
				return
			}
			rentedFromInventory := 0
			for _, rental := range rentals {
				if !rental.ReturnDate.Valid {
					rentedFromInventory++
				}
			}
			movie.AvailableToRent = len(inventories) - rentedFromInventory
		}
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, moviesResponse{Data: movies}); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}
}
