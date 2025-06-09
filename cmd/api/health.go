package main

import (
	"net/http"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type healthCheckData struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type healthCheckResponse struct {
	Data healthCheckData `json:"data"`
}

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the server is running
//	@Tags			1. Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} healthCheckResponse
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := healthCheckData{
		Status:      "ok",
		Environment: app.config.env,
		Version:     version,
	}
	if err := utils.WriteJSONResponse(w, http.StatusOK, healthCheckResponse{Data: data}); err != nil {
		app.errorHandler.InternalServerError(w, r, err)
	}

}
