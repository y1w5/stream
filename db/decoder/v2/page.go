package decoder

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Page represents a page from Wikipedia.
type Page struct {
	ID        int64     `xml:"-"`
	UpdatedAt time.Time `xml:"revision>timestamp"`
	Title     string    `xml:"title"`
	Text      string    `xml:"revision>text"`
}

// UnmarshalXML unmarshals an XML element into the page.
func (p *Page) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.RawToken()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.StartElement:
			if err := p.unmarshalXMLField(d, t); err != nil {
				return err
			}
		case xml.EndElement:
			// We reached the end
			if t.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
}

func (p *Page) unmarshalXMLField(d *xml.Decoder, start xml.StartElement) error {
	field := start.Name.Local
	switch field {
	case "revision":
		var r revisionV3
		if err := r.UnmarshalXML(d, start); err != nil {
			return fmt.Errorf("unmarshal %s: %v", field, err)
		}
		p.UpdatedAt = r.UpdatedAt
		p.Text = r.Text
	case "title":
		s, err := parse(d, start)
		if err != nil {
			return fmt.Errorf("unmarshal %s: %v", field, err)
		}
		p.Title = s
	default:
		return consume(d, start)
	}

	return nil

}

type revisionV3 struct {
	UpdatedAt time.Time `xml:"timestamp"`
	Text      string    `xml:"text"`
}

func (r *revisionV3) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.RawToken()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.StartElement:
			if err := r.unmarshalXMLField(d, t); err != nil {
				return err
			}
		case xml.EndElement:
			// We reached the end
			if t.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
}

func (r *revisionV3) unmarshalXMLField(d *xml.Decoder, start xml.StartElement) error {
	field := start.Name.Local
	switch field {
	case "timestamp":
		s, err := parse(d, start)
		if err != nil {
			return fmt.Errorf("unmarshal %s: %v", field, err)
		}
		r.UpdatedAt, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("unmarshal %s: %v", field, err)
		}
	case "text":
		s, err := parse(d, start)
		if err != nil {
			return fmt.Errorf("unmarshal %s: %v", field, err)
		}
		r.Text = s
	default:
		return consume(d, start)
	}

	return nil

}
