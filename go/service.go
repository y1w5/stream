package main

import (
	"context"
	"fmt"
	"time"
)

// Service stores the business logic of our application.
type Service struct {
	db *DB
}

// NewService instanciates a [Service].
func NewService(dbPath string) (*Service, error) {
	db, err := NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("new db: %v", err)
	}

	return &Service{
		db: db,
	}, nil
}

// Close closes allocated ressources.
func (s *Service) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db: %v", err)
	}
	return nil
}

// Page stores information on a Wiki page.
type Page struct {
	ID        int64
	UpdatedAt time.Time
	Title     string
}

// ListPages lists all pages.
func (s *Service) ListPages(ctx context.Context) ([]Page, error) {
	pages, err := s.db.ListPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("db: %v", err)
	}

	// TODO: apply the business logic
	return pages, nil
}
