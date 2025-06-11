package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCustomerRepository(t *testing.T) {
	ctx := context.Background()
	// Create a temporary directory to store the .dat files and restore.sql
	tempDir, err := os.MkdirTemp("", "dvdrental-data")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Copy all .dat files to the temporary directory
	datFiles, err := filepath.Glob(filepath.Join("..", "..", "testdata", "dvdrental", "*.dat"))
	if err != nil {
		t.Fatal(err)
	}

	for _, datFile := range datFiles {
		content, err := os.ReadFile(datFile)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(filepath.Join(tempDir, filepath.Base(datFile)), content, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Copy restore.sql to the temporary directory
	restoreSQL, err := os.ReadFile(filepath.Join("..", "..", "testdata", "dvdrental", "new_restore.sql"))
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "01-restore.sql"), restoreSQL, 0644)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)

	rentalStore := NewRentalStore(db)

	rental, err := rentalStore.GetRental(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 1, rental.ID)
	assert.Equal(t, "2005-05-24 22:53:30 +0000 +0000", rental.RentalDate.String())

	rental, err = rentalStore.GetRental(ctx, 2)
	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 2, rental.ID)
	assert.Equal(t, "2005-05-24 22:54:33 +0000 +0000", rental.RentalDate.String())
}
