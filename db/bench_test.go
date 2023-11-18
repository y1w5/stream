package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/y1w5/stream/db/decoder"
	decoderv2 "github.com/y1w5/stream/db/decoder/v2"
)

var loadDatasetInMemory = sync.OnceValues(func() ([]byte, error) {
	var err error
	dataset, err := loadDataset()
	if err != nil {
		return nil, fmt.Errorf("fail to load dataset: %w", err)
	}

	buf, err := io.ReadAll(dataset)
	if err != nil {
		return nil, fmt.Errorf("fail to read buffer: %v", err)
	}

	return buf, nil
})

func BenchmarkDecoder(b *testing.B) {
	dataset, err := loadDatasetInMemory()
	if errors.Is(err, ErrDatasetNotFound) {
		b.SkipNow()
	}
	if err != nil {
		b.Fatal(err)
	}

	d, err := decoder.New(bytes.NewReader(dataset))
	if err != nil {
		b.Fatalf("fail to create decoder: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p decoder.Page

		if !d.Next() {
			b.Fatalf("fail to read next page: %v", err)
		}

		err := d.Scan(&p)
		if err != nil {
			b.Fatalf("fail to scan page: %v", err)
		}
	}
}

func BenchmarkDecoderV2(b *testing.B) {
	dataset, err := loadDatasetInMemory()
	if errors.Is(err, ErrDatasetNotFound) {
		b.SkipNow()
	}
	if err != nil {
		b.Fatal(err)
	}

	decoder, err := decoderv2.New(bytes.NewReader(dataset))
	if err != nil {
		b.Fatalf("fail to create decoder: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p decoderv2.Page

		if !decoder.Next() {
			if err != nil {
				b.Fatalf("fail to read next page: %v", err)
			}
		}

		err := decoder.Scan(&p)
		if err != nil {
			b.Fatalf("fail to scan page: %v", err)
		}
	}
}

func BenchmarkDecoder_streaming(b *testing.B) {
	dataset, err := loadDataset()
	if errors.Is(err, ErrDatasetNotFound) {
		b.SkipNow()
	}
	if err != nil {
		b.Fatal(err)
	}

	d, err := decoder.New(dataset)
	if err != nil {
		b.Fatalf("fail to create decoder: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p decoder.Page

		if !d.Next() {
			b.Fatalf("fail to read next page: %v", err)
		}

		err := d.Scan(&p)
		if err != nil {
			b.Fatalf("fail to scan page: %v", err)
		}
	}
}

func BenchmarkSummarize(b *testing.B) {
	dataset, err := loadDatasetInMemory()
	if errors.Is(err, ErrDatasetNotFound) {
		b.SkipNow()
	}
	if err != nil {
		b.Fatal(err)
	}

	d, err := decoder.New(bytes.NewReader(dataset))
	if err != nil {
		b.Fatalf("fail to create decoder: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p decoder.Page

		if !d.Next() {
			b.Fatalf("fail to read next page: %v", err)
		}

		err := d.Scan(&p)
		if err != nil {
			b.Fatalf("fail to scan page: %v", err)
		}

		Summarize(p.Text)
	}
}
