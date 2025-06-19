package main

import (
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

// GetMovies godoc
//
//	@Summary		Get movies
//	@Description	Get all movies
//	@Tags			5. Movies
//	@Accept			json
//	@Produce		json
//	@Param			query	query		utils.MovieQuery	true	"Query parameters"
//	@Success		200		{object}	any
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

	if err := utils.WriteJSONResponse(w, http.StatusOK, nil); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
		return
	}
}
