package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/cheggaaa/pb/v3"
	"github.com/y1w5/stream/db/decoder/v2"
)

var spinner pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}} {{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" }} {{counters .}} {{speed . "%s p/s"}} {{with string . "suffix"}} {{.}}{{end}}`
var spinnerETA pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}} {{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" }} {{counters .}} {{speed . "%s p/s"}} {{rtime .}}{{with string . "suffix"}} {{.}}{{end}}`

func main() {
	fmt.Printf("Loading Wikipedia dataset...\n")
	datasets, err := loadDatasets()
	if err != nil {
		fatalf("fail to load datasetis: %v", err)
	}
	defer datasets.Close()

	fmt.Printf("Setting up SQLite database...\n")
	db, err := newDB()
	if err != nil {
		fatalf("fail to create db: %v", err)
	}
	defer db.Close()

	fmt.Printf("Loading dataset into SQLite...\n")
	count := 0
	for _, d := range datasets {
		bar := spinner.Start(0).Set("prefix", "  "+d.Name())
		c, err := insertDataset(bar, db, d)
		if err != nil {
			fatalf("fail to insert dataset: %v", err)
		}
		count += c
		bar.Finish()
	}

	fmt.Printf("Completed, %d pages created.\n", count)
}

func loadDatasets() (Datasets, error) {
	type result struct {
		Dataset *Dataset
		Err     error
	}

	pool := pb.NewPool()
	results := make(chan result)
	for _, name := range datasets {
		bar := spinnerETA.New(0).
			Set(pb.Bytes, true).
			Set("prefix", "  "+name)
		pool.Add(bar)

		go func(name string, bar *pb.ProgressBar) {
			defer bar.Finish()

			d, err := loadDataset(name)
			if errors.Is(err, ErrDatasetNotFound) {
				d, err = downloadDataset(bar, name)
				results <- result{d, err}
				return
			}
			if err != nil {
				results <- result{nil, err}
				return
			}

			bar.AddTotal(d.Size())
			bar.Add64(d.Size())
			results <- result{d, nil}
		}(name, bar)
	}

	err := pool.Start()
	if err != nil {
		return nil, fmt.Errorf("start pool: %v", err)
	}
	defer pool.Stop()

	datasets := make([]*Dataset, len(datasets))
	for i := range datasets {
		result := <-results
		if result.Err != nil {
			return nil, result.Err
		}
		datasets[i] = result.Dataset
	}

	slices.SortFunc(datasets, func(a, b *Dataset) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	return datasets, nil
}

func insertDataset(bar *pb.ProgressBar, db *DB, dataset *Dataset) (int, error) {
	d, err := decoder.New(dataset)
	if err != nil {
		return 0, fmt.Errorf("create decoder: %v\n", err)
	}

	if err := db.Begin(); err != nil {
		return 0, fmt.Errorf("begin transaction: %v\n", err)
	}
	defer db.Rollback() //nolint

	var count int
	for d.Next() {
		var p decoder.Page

		err := d.Scan(&p)
		if err != nil {
			return 0, fmt.Errorf("scan page: %v\n", err)
		}

		_, err = db.CreatePage(CreatePageParams{
			UpdatedAt: p.UpdatedAt,
			Title:     p.Title,
			Text:      Summarize(p.Text),
		})
		if err != nil {
			return 0, fmt.Errorf("create page: %v\n", err)
		}

		count++
		bar.Increment()
	}
	if err := d.Err(); err != nil && !errors.Is(err, io.EOF) {
		return 0, fmt.Errorf("scan pages: %v\n", err)
	}
	if err := db.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %v\n", err)
	}
	return count, nil
}

func fatalf(format string, v ...any) {
	fmt.Printf(format, v...)
	os.Exit(1)
}
