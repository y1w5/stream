package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/y1w5/stream/db/decoder/v2"
)

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
	defer db.Rollback() //nolint

	fmt.Printf("Loading dataset into SQLite...\n")
	d, err := decoder.New(dataset)
	if err != nil {
		fatalf("fail to create decoder: %v\n", err)
	}

	var count int
	for d.Next() {
		var p decoder.Page

		err := d.Scan(&p)
		if err != nil {
			fatalf("fail to scan page: %v\n", err)
		}

		_, err = db.CreatePage(CreatePageParams{
			UpdatedAt: p.UpdatedAt,
			Title:     p.Title,
			Text:      Summarize(p.Text),
		})
		if err != nil {
			fatalf("fail to create page: %v\n", err)
		}

		count++
		if count%1000 == 0 {
			fmt.Printf("Copied %d pages.\n", count)
		}
	}
	if err := d.Err(); err != nil && !errors.Is(err, io.EOF) {
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
