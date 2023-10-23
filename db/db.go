package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const (
	dbName   = "./stream.db"
	dbSchema = "./schema.sql"
)

type dbtx interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

type DB struct {
	db   *sql.DB
	tx   *sql.Tx
	dbtx dbtx
}

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

type CreatePageParams struct {
	UpdatedAt time.Time
	Title     string
	Text      string
}

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

func (db *DB) Commit() error {
	if db.tx == nil {
		return fmt.Errorf("Commit called outside transaction")
	}
	tx := db.tx
	db.tx, db.dbtx = nil, db.db
	return tx.Commit()
}

func (db *DB) Rollback() error {
	if db.tx == nil {
		return fmt.Errorf("Rollback called outside transaction")
	}
	tx := db.tx
	db.tx, db.dbtx = nil, db.db
	return tx.Rollback()
}

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
