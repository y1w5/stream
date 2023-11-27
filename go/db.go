package main

import (
	"context"
	"database/sql"
	"fmt"
)

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

var listPagesQuery = `SELECT id, updated_at, title, text FROM pages`

// ListPages lists all pages.
func (db *DB) ListPages(ctx context.Context) ([]Page, error) {
	rows, err := db.db.QueryContext(ctx, listPagesQuery)
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
func (db *DB) StreamPages(ctx context.Context) func(func(Page, error) bool) {
	return func(yield func(Page, error) bool) {
		var zero Page
		rows, err := db.db.QueryContext(ctx, listPagesQuery)
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

		return
	}
}
