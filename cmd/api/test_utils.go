package main

import (
	"testing"

	"github.com/andras-szesztai/dev-rental-api/internal/auth"
	"github.com/andras-szesztai/dev-rental-api/internal/cache"
	"github.com/andras-szesztai/dev-rental-api/internal/store"
	"github.com/andras-szesztai/dev-rental-api/internal/utils"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	return &application{
		logger:        zap.NewNop().Sugar(),
		store:         store.NewMockStore(),
		authenticator: auth.NewMockAuth(),
		errorHandler:  utils.NewErrorHandler(zap.NewNop().Sugar()),
		cache:         cache.NewMockCache(),
	}
}
