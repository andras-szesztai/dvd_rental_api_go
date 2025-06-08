package main

import (
	"log"
	"os"

	"go.uber.org/zap"

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
	logger *zap.SugaredLogger
	config config
	store  *store.Store
}

type config struct {
	addr string
	env  string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	app := &application{
		logger: logger,
		config: cfg,
		store:  store,
	}

	err = app.serve(app.mountRoutes())
	if err != nil {
		logger.Fatalw("failed to serve", "error", err)
	}
}
