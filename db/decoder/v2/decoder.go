// Package decoder implements a high performance decoder for wiki pages.
//
// It uses xml.Decoder.RawToken to speed up parsing and Page.UnmarshalXML
// to reduce allocation.
package decoder

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Decoder is an XML decoder tailored to the Wikipedia dataset.
type Decoder struct {
	err   error
	start xml.StartElement

	d *xml.Decoder
}

// New instanciates a new Decoder.
//
// It fail if it cannot find the mediawiki and siteinfo elements from the dataset.
func New(r io.Reader) (*Decoder, error) {
	// Instanciates the decoder.
	d := &Decoder{
		d: xml.NewDecoder(r),
	}

	// Read the mediawiki element.
	err := d.checkToken("mediawiki")
	if err != nil {
		return nil, fmt.Errorf("check mediawiki: %v", err)
	}

	// Read the siteinfo element.
	err = d.checkToken("siteinfo")
	if err != nil {
		return nil, fmt.Errorf("check siteinfo: %v", err)
	}
	if !d.consume() {
		return nil, d.err
	}

	return d, nil
}

// Next moves to the next element.
func (d *Decoder) Next() bool {
	if d.err != nil {
		return false
	}

	return d.next()
}

// Err returns any error encountered by Next.
func (d *Decoder) Err() error { return d.err }

// Scan scans an element.
//
// Beware, UnmarshalXML must be implemented by calling RawToken.
func (d *Decoder) Scan(v xml.Unmarshaler) error {
	if d.err != nil {
		return d.err
	}
	return v.UnmarshalXML(d.d, d.start)
}

func (d *Decoder) checkToken(name string) error {
	if !d.next() {
		return d.err
	}

	n := d.start.Name.Local
	if n != name {
		return fmt.Errorf("unexpected name: expects=%v got=%v", n, name)
	}
	return nil
}

// consume consumes the current element and move to the next one.
func (d *Decoder) consume() bool {
	var zero xml.StartElement

	err := consume(d.d, d.start)
	if err != nil {
		d.start, d.err = zero, err
		return false
	}

	return true
}

// next returns the next element without consuming the current element.
func (d *Decoder) next() bool {
	var zero xml.StartElement

	start, err := next(d.d)
	if err != nil {
		d.start, d.err = zero, err
		return false
	}

	d.start = start
	return true
}
