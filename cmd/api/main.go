package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	"github.com/andras-szesztai/dev-rental-api/internal/auth"
	"github.com/andras-szesztai/dev-rental-api/internal/db"
	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"github.com/joho/godotenv"
)

func getVersion() string {
	version, err := os.ReadFile("version.txt")
	if err != nil {
		return "unknown"
	}
	return string(version)
}

//	@title			DVD Rental API
//	@version		0.0.1
//	@description	API for a DVD Rental management application
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//	@scheme						bearer
//	@type						http

type application struct {
	logger        *zap.SugaredLogger
	config        config
	store         *store.Store
	authenticator auth.Authenticator
	errorHandler  *utils.ErrorHandler
}

type config struct {
	addr    string
	env     string
	db      dbConfig
	auth    authConfig
	apiURL  string
	version string
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

//	@title			Swagger Examasdasdasdasdasdawdasple API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//	@scheme						bearer
//	@type						http

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	exp, err := time.ParseDuration(os.Getenv("TOKEN_EXP"))
	if err != nil {
		fmt.Println("Error parsing TOKEN_EXP", err)
	}

	cfg := config{
		addr:    os.Getenv("PORT_ADDR"),
		env:     os.Getenv("ENV"),
		version: getVersion(),
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
		apiURL: os.Getenv("API_URL"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal("failed to sync logger", "error", err)
		}
	}()

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStore(db)

	authenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	errorHandler := utils.NewErrorHandler(logger)

	app := &application{
		logger:        logger,
		config:        cfg,
		store:         store,
		authenticator: authenticator,
		errorHandler:  errorHandler,
	}

	err = app.serve(app.mountRoutes())
	if err != nil {
		logger.Fatalw("failed to serve", "error", err)
	}
}
