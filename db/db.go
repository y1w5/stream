package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName   = "./stream.db"
	dbSchema = "./schema.sql"
)

type dbtx interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Page represents a page in the database.
type Page struct {
	ID        int64
	UpdatedAt time.Time
	Title     string
	Text      string
}

// DB represents the database access layer of our application.
type DB struct {
	db   *sql.DB
	tx   *sql.Tx
	dbtx dbtx
}

// newDB instanciates a new database.
//
// It drops existing SQLite database and recreate a new one from scratch.
func newDB() (*DB, error) {
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", dbName, err)
	}

	err = migrateUp(db)
	if err != nil {
		return nil, fmt.Errorf("migrate up: %v", err)
	}

	return &DB{
		db:   db,
		dbtx: db,
	}, nil
}

var createPageQuery = `INSERT INTO pages (updated_at, title, "text")
VALUES (?, ?, ?)
RETURNING id, updated_at, title, text;`

// CreatePageParams stores required parameters for [CreatePage].
type CreatePageParams struct {
	UpdatedAt time.Time
	Title     string
	Text      string
}

// CreatePage creates a new page in the database.
func (db *DB) CreatePage(arg CreatePageParams) (Page, error) {
	var p Page

	err := db.dbtx.QueryRow(createPageQuery,
		arg.UpdatedAt.Format(time.DateTime),
		arg.Title,
		arg.Text,
	).Scan(&p.ID, &p.UpdatedAt, &p.Title, &p.Text)
	if err != nil {
		return Page{}, err
	}

	return p, nil
}

// Begin begins a transaction.
func (db *DB) Begin() error {
	if db.tx != nil {
		return fmt.Errorf("a transaction is already running")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	db.tx, db.dbtx = tx, tx
	return nil
}

// Commit commits a transaction.
func (db *DB) Commit() error {
	if db.tx == nil {
		return fmt.Errorf("Commit called outside transaction")
	}
	tx := db.tx
	db.tx, db.dbtx = nil, db.db
	return tx.Commit()
}

// Rollback rollbacks a transaction.
func (db *DB) Rollback() error {
	if db.tx == nil {
		return fmt.Errorf("Rollback called outside transaction")
	}
	tx := db.tx
	db.tx, db.dbtx = nil, db.db
	return tx.Rollback()
}

// Close closes the underling database.
func (db *DB) Close() error {
	if db.tx != nil {
		_ = db.Rollback()
	}
	return db.db.Close()
}

func migrateUp(db *sql.DB) error {
	schema, err := os.ReadFile(dbSchema)
	if err != nil {
		return fmt.Errorf("read file %v: %v", dbSchema, err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
