package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	ID        int64     `xml:"-"`
	UpdatedAt time.Time `xml:"revision>timestamp"`
	Title     string    `xml:"title"`
	Text      string    `xml:"revision>text"`
}

func main() {
	fmt.Printf("Loading Wikipedia dataset...\n")
	dataset, err := loadDataset()
	if errors.Is(err, ErrDatasetNotFound) {
		dataset, err = downloadDataset()
	}
	if err != nil {
		fatalf("fail to load dataset: %v", err)
	}
	defer dataset.Close()

	fmt.Printf("Setting up SQLite database...\n")
	db, err := newDB()
	if err != nil {
		fatalf("fail to create db: %v", err)
	}
	defer db.Close()

	if err := db.Begin(); err != nil {
		fatalf("fail to begin transaction: %v\n", err)
	}
	defer db.Rollback()

	fmt.Printf("Loading dataset into SQLite...\n")
	decoder, err := newDecoder(dataset)
	if err != nil {
		fatalf("fail to create decoder: %v\n", err)
	}

	var count int
	for decoder.Next() {
		var p Page

		err := decoder.Scan(&p)
		if err != nil {
			fatalf("fail to scan page: %v\n", err)
		}

		_, err = db.CreatePage(CreatePageParams{
			UpdatedAt: p.UpdatedAt,
			Title:     p.Title,
			Text:      p.Text,
		})
		if err != nil {
			fatalf("fail to create page: %v\n", err)
		}

		count++
		if count%1000 == 0 {
			fmt.Printf("Copied %d pages.\n", count)
		}
	}
	if err := decoder.Err(); err != nil && !errors.Is(err, io.EOF) {
		fatalf("fail to scan pages: %v\n", err)
	}
	if err := db.Commit(); err != nil {
		fatalf("fail to commit transaction: %v\n", err)
	}

	fmt.Printf("Completed, %d pages created.\n", count)
}

func fatalf(format string, v ...any) {
	fmt.Printf(format, v...)
	os.Exit(1)
}
