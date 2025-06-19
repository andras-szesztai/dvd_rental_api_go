package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andras-szesztai/dev-rental-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type contextKey string

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func (app *application) mountRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(app.RateLimiterMiddleware)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/test", app.testRateLimiter)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
			httpSwagger.DomID("swagger-ui"),
			httpSwagger.UIConfig(map[string]string{
				"tagsSorter":       "\"alpha\"",
				"operationsSorter": "\"method\"",
			}),
		))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUser)
			r.Post("/sign-in", app.signInUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Route("/rentals", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.getRentalByID)
				})
			})
			r.Route("/customers", func(r chi.Router) {
				r.Post("/", app.CheckAdminMiddleware(app.createCustomer))
			})
			r.Route("/movies", func(r chi.Router) {
				r.Get("/", app.getMovies)
			})
		})
	})

	return router
}

func (app *application) serve(router http.Handler) error {
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Title = "DVD Rental API"
	docs.SwaggerInfo.Description = "API for a DVD Rental management application"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("starting server", "addr", srv.Addr, "env", app.config.env, "version", app.config.version)

	return srv.ListenAndServe()
}

func (app *application) testRateLimiter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Rate limiter test endpoint", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}
