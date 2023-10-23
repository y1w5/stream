package main

import (
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	datasetName = "enwiki-20231020-pages-articles1.xml-p1p41242.bz2"
	datasetURL  = "https://dumps.wikimedia.org/enwiki/20231020/enwiki-20231020-pages-articles1.xml-p1p41242.bz2"
)

var ErrDatasetNotFound = errors.New("dataset: not found")

type Dataset struct {
	f *os.File
	io.Reader
}

func loadDataset() (*Dataset, error) {
	f, err := os.Open(datasetName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrDatasetNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", datasetName, err)
	}

	return &Dataset{
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

func downloadDataset() (*Dataset, error) {

	resp, err := http.Get(datasetURL)
	if err != nil {
		return nil, fmt.Errorf("get %v: %v", datasetURL, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(datasetName)
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", datasetName, err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return nil, fmt.Errorf("copy: %v", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("seek: %v", err)
	}

	return &Dataset{
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

func (d *Dataset) Close() error { return d.f.Close() }
