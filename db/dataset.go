package main

import (
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const datasetExt = ".bz2"

const datasetBaseURL = "https://dumps.wikimedia.org/enwiki/20231020/"

var datasets = []string{
	"enwiki-20231020-pages-articles1.xml-p1p41242",
	"enwiki-20231020-pages-articles2.xml-p41243p151573",
	"enwiki-20231020-pages-articles3.xml-p151574p311329",
	"enwiki-20231020-pages-articles4.xml-p311330p558391",
	"enwiki-20231020-pages-articles5.xml-p558392p958045",
	"enwiki-20231020-pages-articles6.xml-p958046p1483661",
}

var ErrDatasetNotFound = errors.New("dataset: not found")

type Dataset struct {
	name string
	size int64
	f    *os.File
	io.Reader
}

func loadDataset(name string) (*Dataset, error) {
	var err error
	var f *os.File
	for _, path := range []string{name, name + datasetExt} {
		f, err = os.Open(path)
		if err == nil {
			break
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		return nil, fmt.Errorf("open %v: %v", path, err)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrDatasetNotFound
	}

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %v", err)
	}

	return &Dataset{
		name:   name,
		size:   info.Size(),
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

func downloadDataset(bar *pb.ProgressBar, name string) (*Dataset, error) {
	url := datasetBaseURL + name + datasetExt

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get %v: %v", url, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(name + datasetExt)
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", name, err)
	}

	reader := bar.AddTotal(resp.ContentLength).NewProxyReader(resp.Body)
	if _, err := io.Copy(f, reader); err != nil {
		return nil, fmt.Errorf("copy: %v", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("seek: %v", err)
	}

	return &Dataset{
		name:   name,
		size:   resp.ContentLength,
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

func (d *Dataset) Name() string { return d.name }

func (d *Dataset) Size() int64 { return d.size }

func (d *Dataset) Close() error {
	f := d.f
	d.f = nil
	d.Reader = nil
	return f.Close()
}

type Datasets []*Dataset

func (datasets Datasets) Close() error {
	var errs []error
	for _, d := range datasets {
		if err := d.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
