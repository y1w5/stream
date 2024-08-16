package main

import (
	"context"
	"testing"
)

const dbPath = "stream.db"

func BenchmarkDBListPages(b *testing.B) {
	db, err := NewDB(dbPath)
	if err != nil {
		b.Fatalf("new DB: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	ctx := context.Background()
	for range b.N {
		pages, err := db.ListPages(ctx, -1)
		if err != nil {
			b.Fatalf("list pages: %v", err)
		}
		_ = pages
	}
}

func BenchmarkDBStreamPages(b *testing.B) {
	db, err := NewDB(dbPath)
	if err != nil {
		b.Fatalf("new DB: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	ctx := context.Background()
	for range b.N {
		for p, err := range db.StreamPages(ctx, -1) {
			if err != nil {
				b.Fatalf("stream pages: %v", err)
			}
			_ = p
		}
	}
}

func BenchmarkDBStreamPageSlice(b *testing.B) {
	db, err := NewDB(dbPath)
	if err != nil {
		b.Fatalf("new DB: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	ctx := context.Background()
	for range b.N {
		for pages, err := range db.StreamPageSlice(ctx, -1) {
			if err != nil {
				b.Fatalf("stream pages slice: %v", err)
			}
			_ = pages
		}
	}
}
