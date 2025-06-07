package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) mountRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/rentals/{id}", app.getRentalByID)
	})

	return router
}

func (app *application) serve(router http.Handler) error {
	// Docs
	// docs.SwaggerInfo.Version = version
	// docs.SwaggerInfo.Host = app.config.apiURL
	// docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("starting server", "addr", srv.Addr, "env", app.config.env, "version", version)

	return srv.ListenAndServe()
}
