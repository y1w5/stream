package main

import (
	"strings"
)

// Sumamrize reduces the given text to a few paragraphs.
//
// It strips curly braces and returns all the text until it hits a title.
func Summarize(text string) string {
	const maxLineBreak = 2
	var lineBreaks = make([]byte, 0, maxLineBreak)

	var w strings.Builder
	r := strings.NewReader(text)
	skipSpaces(r)
loop:
	for {
		c, err := r.ReadByte()
		if err != nil {
			break
		}
		switch {
		case c == '\n':
			if len(lineBreaks) < maxLineBreak && w.Len() > 0 {
				lineBreaks = append(lineBreaks, c)
			}
		case c == '{' && consumeByte(r, '{'):
			skipCurlyBraces(r)
			continue
		case c == '=' && consumeByte(r, '='):
			break loop
		default:
			if len(lineBreaks) > 0 {
				_, _ = w.Write(lineBreaks)
				lineBreaks = lineBreaks[:0]
			}
			_ = w.WriteByte(c)
		}
	}

	return w.String()
}

func consumeByte(r *strings.Reader, c byte) bool {
	cc, err := r.ReadByte()
	if err != nil {
		return false
	}
	if c == cc {
		return true
	}
	_ = r.UnreadByte()
	return false
}

func skipCurlyBraces(r *strings.Reader) {
	count := 1
	for count != 0 {
		c, err := r.ReadByte()
		if err != nil {
			break
		}

		if c == '{' && consumeByte(r, '{') {
			count++
		}
		if c == '}' && consumeByte(r, '}') {
			count--
		}
	}
}

func skipSpaces(r *strings.Reader) {
	for {
		c, err := r.ReadByte()
		if err != nil {
			break
		}

		if c != ' ' && c != '\t' {
			_ = r.UnreadByte()
			break
		}
	}
}
