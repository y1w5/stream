package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DBSliceSize is the size of a slice of data.
//
// It was selected by trial and error to find the best performing buffer size.
const DBSliceSize = 1 << 16

// DB is the database access layer of our application.
type DB struct {
	db *sql.DB
}

// NewDB instanciates a [DB].
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", path, err)
	}

	return &DB{
		db: db,
	}, nil
}

// Close closes allocated ressources.
func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return err
	}
	return nil
}

var listPagesQuery = `SELECT id, updated_at, title, text FROM pages LIMIT ?`

// Page stores information on a Wiki page.
type Page struct {
	ID        int64
	UpdatedAt time.Time
	Title     string
	Text      string
}

// ListPages lists all pages.
func (db *DB) ListPages(ctx context.Context, limit int) ([]Page, error) {
	rows, err := db.db.QueryContext(ctx, listPagesQuery, softLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer rows.Close()

	var pages []Page
	for rows.Next() {
		var p Page
		err := rows.Scan(&p.ID, &p.UpdatedAt, &p.Title, &p.Text)
		if err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}
		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next: %v", err)
	}

	return pages, nil
}

// StreamPages streams pages from the database.
func (db *DB) StreamPages(ctx context.Context, limit int) func(func(Page, error) bool) {
	return func(yield func(Page, error) bool) {
		var zero Page
		rows, err := db.db.QueryContext(ctx, listPagesQuery, softLimit(limit))
		if err != nil {
			yield(zero, fmt.Errorf("query: %v", err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var p Page
			err := rows.Scan(&p.ID, &p.UpdatedAt, &p.Title, &p.Text)
			if err != nil {
				yield(zero, fmt.Errorf("scan: %v", err))
				return
			}
			if !yield(p, err) {
				return
			}
		}

		if err := rows.Err(); err != nil {
			yield(zero, fmt.Errorf("next: %v", err))
			return
		}
	}
}

// StreamPageSlice streams pages from the database into slices.
func (db *DB) StreamPageSlice(ctx context.Context, limit int) func(func([]Page, error) bool) {
	return func(yield func([]Page, error) bool) {
		rows, err := db.db.QueryContext(ctx, listPagesQuery, softLimit(limit))
		if err != nil {
			yield(nil, fmt.Errorf("query: %v", err))
			return
		}
		defer rows.Close()

		pages := make([]Page, 0, DBSliceSize)
		for rows.Next() {
			var p Page
			err := rows.Scan(&p.ID, &p.UpdatedAt, &p.Title, &p.Text)
			if err != nil {
				yield(nil, fmt.Errorf("scan: %v", err))
				return
			}

			pages = append(pages, p)
			if len(pages) < DBSliceSize {
				continue
			}

			if !yield(pages, err) {
				return
			}
			pages = pages[:0]
		}

		if err := rows.Err(); err != nil {
			yield(nil, fmt.Errorf("next: %v", err))
			return
		}

		if len(pages) > 0 {
			yield(pages, nil)
		}
	}
}

// softLimit changes the zero value of limit into -1, allowing SQLite to return
// the full dataset if the limit is unset.
func softLimit(limit int) int {
	if limit < 1 {
		return -1
	}
	return limit
}
