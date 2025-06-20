package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/andras-szesztai/dev-rental-api/internal/utils"
)

type MovieStore struct {
	db *sql.DB
}

func NewMovieStore(db *sql.DB) *MovieStore {
	return &MovieStore{db: db}
}

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
}

func (s *MovieStore) GetMovies(ctx context.Context, movieQuery *utils.MovieQuery) ([]*Movie, error) {
	queryStr := `
		SELECT film_id, title, description, release_year	
		FROM film
		WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
		ORDER BY release_year ` + movieQuery.Sort + ` 
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, queryStr, movieQuery.Search, movieQuery.Limit, movieQuery.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
