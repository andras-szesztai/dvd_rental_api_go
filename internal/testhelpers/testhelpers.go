package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
	DB               *sql.DB
}

func CreatePostgresContainer() (*PostgresContainer, error) {
	ctx := context.Background()
	// Create a temporary directory to store the .dat files and restore.sql
	tempDir, err := os.MkdirTemp("", "dvdrental-data")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Copy all .dat files to the temporary directory
	datFiles, err := filepath.Glob(filepath.Join("..", "..", "testdata", "dvdrental", "*.dat"))
	if err != nil {
		return nil, err
	}

	for _, datFile := range datFiles {
		content, err := os.ReadFile(datFile)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(filepath.Join(tempDir, filepath.Base(datFile)), content, 0644)
		if err != nil {
			return nil, err
		}
	}

	// Copy restore.sql to the temporary directory
	restoreSQL, err := os.ReadFile(filepath.Join("..", "..", "testdata", "dvdrental", "new_restore.sql"))
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filepath.Join(tempDir, "01-restore.sql"), restoreSQL, 0644)
	if err != nil {
		return nil, err
	}

	// Copy migration files to the temporary directory
	migrationFiles, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil {
		return nil, err
	}

	for i, migrationFile := range migrationFiles {
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			return nil, err
		}

		// Use a higher number prefix to ensure migrations run after the restore
		err = os.WriteFile(filepath.Join(tempDir, fmt.Sprintf("02-%03d-migration.sql", i+1)), content, 0644)
		if err != nil {
			return nil, err
		}
	}

	pgContainer, err := postgres.Run(ctx,
		"postgres:16.1-alpine",
		postgres.WithDatabase("dvdrental"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(10*time.Second)),
		testcontainers.WithMounts(testcontainers.ContainerMount{
			Source:   testcontainers.GenericBindMountSource{HostPath: tempDir},
			Target:   "/docker-entrypoint-initdb.d",
			ReadOnly: true,
		}),
	)

	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
		DB:                db,
	}, nil
}
