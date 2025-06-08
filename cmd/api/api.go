package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

func (app *application) mountRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUser)
			r.Post("/sign-in", app.signInUser)
		})
		// Auth routes
		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Route("/rentals", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.getRentalByID)
				})
			})

		})
	})

	return router
}

func (app *application) serve(router http.Handler) error {
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
