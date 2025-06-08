package main

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/andras-szesztai/dev-rental-api/internal/auth"
	"github.com/andras-szesztai/dev-rental-api/internal/db"
	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/joho/godotenv"
)

// 1. Setup swagger
// 2. Setup POST rental (check inventory)
// 4. Setup PATCH rental for return
// 5. Setup AUTH (JWT)
// 6. Setup Authorization

const version = "0.0.1"

type application struct {
	logger        *zap.SugaredLogger
	config        config
	store         *store.Store
	authenticator auth.Authenticator
}

type config struct {
	addr string
	env  string
	db   dbConfig
	auth authConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	aud    string
	iss    string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	exp, err := time.ParseDuration(os.Getenv("TOKEN_EXP"))
	if err != nil {
		log.Fatal("Error parsing TOKEN_EXP")
	}

	cfg := config{
		addr: os.Getenv("PORT"),
		env:  os.Getenv("ENV"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: 50,
			maxIdleConns: 25,
			maxIdleTime:  "15m",
		},
		auth: authConfig{
			token: tokenConfig{
				secret: os.Getenv("TOKEN_SECRET"),
				exp:    exp,
				aud:    os.Getenv("TOKEN_AUD"),
				iss:    os.Getenv("TOKEN_ISS"),
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStore(db)

	authenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	app := &application{
		logger:        logger,
		config:        cfg,
		store:         store,
		authenticator: authenticator,
	}

	err = app.serve(app.mountRoutes())
	if err != nil {
		logger.Fatalw("failed to serve", "error", err)
	}
}
