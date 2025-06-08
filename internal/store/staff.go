package store

import (
	"context"
	"database/sql"
)

type StaffStore struct {
	db *sql.DB
}

func NewStaffStore(db *sql.DB) *StaffStore {
	return &StaffStore{db: db}
}

type Staff struct {
	ID     int  `json:"id"`
	UserID *int `json:"user_id"`
}

func (s *StaffStore) GetStaffByEmail(ctx context.Context, email string) (*Staff, error) {
	query := `
		SELECT staff_id, user_id
		FROM staff
		WHERE email = $1
	`

	row := s.db.QueryRowContext(ctx, query, email)

	var staff Staff
	var userID sql.NullInt64
	err := row.Scan(&staff.ID, &userID)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		userIDInt := int(userID.Int64)
		staff.UserID = &userIDInt
	}

	return &staff, nil
}
