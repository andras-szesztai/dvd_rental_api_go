package store

import (
	"context"
	"database/sql"
	"time"
)

type InventoryStore struct {
	db *sql.DB
}

func NewInventoryStore(db *sql.DB) *InventoryStore {
	return &InventoryStore{db: db}
}

type Inventory struct {
	ID     int `json:"id"`
	FilmID int `json:"film_id"`
}

func (s *InventoryStore) GetMovieInventory(ctx context.Context, filmID int) ([]*Inventory, error) {
	query := `
		SELECT inventory_id, film_id
		FROM inventory
		WHERE film_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventories := []*Inventory{}
	for rows.Next() {
		var inventory Inventory
		if err := rows.Scan(&inventory.ID, &inventory.FilmID); err != nil {
			return nil, err
		}
		inventories = append(inventories, &inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inventories, nil
}
