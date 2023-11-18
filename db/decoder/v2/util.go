package decoder

import "encoding/xml"

func consume(d *xml.Decoder, start xml.StartElement) error {
	count := 1
	name := start.Name.Local
	for count > 0 {
		t, err := d.RawToken()
		if err != nil {
			return err
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == name {
				count++
			}
		case xml.EndElement:
			if t.Name.Local == name {
				count--
			}
		}
	}

	return nil
}

// next returns the next element without consuming the current element.
func next(d *xml.Decoder) (xml.StartElement, error) {
	var zero xml.StartElement
	for {
		t, err := d.RawToken()
		if err != nil {
			return zero, err
		}

		start, ok := t.(xml.StartElement)
		if ok {
			return start, nil
		}
	}
}

func parse(d *xml.Decoder, start xml.StartElement) (string, error) {
	var str string

	count := 1
	name := start.Name.Local
	for count > 0 {
		t, err := d.RawToken()
		if err != nil {
			return "", err
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == name {
				count++
			}
		case xml.EndElement:
			if t.Name.Local == name {
				count--
			}
		case xml.CharData:
			str = string(t)
		}
	}

	return str, nil
}
