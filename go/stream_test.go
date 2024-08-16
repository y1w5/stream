package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

const streamStdBodyLen = 578_257_270
const streamExpBodyLen = 557_377_380
const streamBufferSize = 600_000_000

func BenchmarkStream(b *testing.B) {
	s, err := NewStream(NewStreamParams{
		Bind:   "localhost:8080",
		DB:     dbPath,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err != nil {
		b.Fatalf("fail to create stream: %v", err)
	}

	benchs := []struct {
		url             string
		method          http.HandlerFunc
		expectedBodyLen int
	}{
		{
			url:             "/test/pages.listStd",
			method:          s.listPagesStd,
			expectedBodyLen: streamStdBodyLen,
		},
		{
			url:             "/test/pages.streamStd",
			method:          s.streamPagesStd,
			expectedBodyLen: 578257268,
		},
		{
			url:             "/test/pages.listExp",
			method:          s.listPagesExp,
			expectedBodyLen: streamExpBodyLen,
		},
		{
			url:             "/test/pages.listSlice",
			method:          s.listPagesSlice,
			expectedBodyLen: streamExpBodyLen,
		},
		{
			url:             "/test/pages.streamSlice",
			method:          s.streamPagesSlice,
			expectedBodyLen: streamExpBodyLen,
		},
		{
			url:             "/pages.list",
			method:          s.listPages,
			expectedBodyLen: streamExpBodyLen,
		},
		{
			url:             "/pages.stream",
			method:          s.streamPages,
			expectedBodyLen: streamExpBodyLen,
		},
		{
			url:             "/pages.streamWithMarshaler",
			method:          s.streamPagesWithMarshaler,
			expectedBodyLen: streamExpBodyLen,
		},
	}

	for _, bb := range benchs {
		b.Run(bb.url[1:], func(b *testing.B) {
			var body bytes.Buffer
			body.Grow(streamBufferSize)

			b.ResetTimer()
			for range b.N {
				resp := &httptest.ResponseRecorder{
					HeaderMap: make(http.Header),
					Body:      &body,
					Code:      200,
				}
				req := httptest.NewRequest("GET", bb.url, nil)
				bb.method(resp, req)
				if resp.Code != 200 {
					b.Fatalf("unexpected status: expects=200 got=%d", resp.Code)
				}
				if resp.Body.Len() != bb.expectedBodyLen {
					b.Fatalf("unexpected body lenght: expects=%d got=%d", bb.expectedBodyLen, resp.Body.Len())
				}
				body.Reset()
			}
		})
	}
}

// listPagesStd lists all pages from the database and write the JSON using
// `encoding/json` from the standard library.
func (s *Stream) listPagesStd(w http.ResponseWriter, r *http.Request) {
	var pages []Page
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		goto encode_err
	}

	pages, err = s.db.ListPages(r.Context(), limit)
	if err != nil {
		s.logger.Error("fail to execute handler", "err", err)
		goto encode_err
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pages)
	return

encode_err:
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
}

// streamPagesStd streams pages from the database and write the JSON using
// `encoding/json` from the standard library.
func (s *Stream) streamPagesStd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	e := json.NewEncoder(w)
	for p, err := range s.db.StreamPages(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		err = e.Encode(p)
		if err != nil {
			s.logger.Error("fail to encode JSON", "err", err)
			return
		}
	}
}

// listPagesExp lists all pages from the database and writes the JSON using
// the experimental `encoding/json/v2`.
func (s *Stream) listPagesExp(w http.ResponseWriter, r *http.Request) {
	var pages []Page
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		goto encode_err
	}

	pages, err = s.db.ListPages(r.Context(), limit)
	if err != nil {
		s.logger.Error("fail to execute handler", "err", err)
		goto encode_err
	}

	w.WriteHeader(http.StatusOK)
	_ = jsonv2.MarshalEncode(jsontext.NewEncoder(w), pages)
	return

encode_err:
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
}

// listPagesSlice lists pages as slice of 65Mo, aggregates them and write the
// JSON using the experimental `encoding/json/v2`.
func (s *Stream) listPagesSlice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	var pages []Page
	for tmps, err := range s.db.StreamPageSlice(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		pages = append(pages, tmps...)
	}

	e := jsontext.NewEncoder(w)
	err = jsonv2.MarshalEncode(e, pages)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

// streamPagesSlice streams pages as slice of 65Mo and write the JSON using
// the experimental `encoding/json/v2`.
func (s *Stream) streamPagesSlice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	e := jsontext.NewEncoder(w)
	err = e.WriteToken(jsontext.ArrayStart)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}

	for pages, err := range s.db.StreamPageSlice(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		for _, p := range pages {
			err = jsonv2.MarshalEncode(e, p)
			if err != nil {
				s.logger.Error("fail to encode JSON", "err", err)
				return
			}
		}
	}

	err = e.WriteToken(jsontext.ArrayEnd)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

func marshalPage(e *jsontext.Encoder, p *Page, opts jsonv2.Options) error {
	_ = e.WriteToken(jsontext.ObjectStart)
	_ = e.WriteToken(jsontext.String("ID"))
	_ = e.WriteToken(jsontext.Int(p.ID))
	_ = e.WriteToken(jsontext.String("UpdatedAt"))
	_ = e.WriteToken(jsontext.String(p.UpdatedAt.Format(time.RFC3339)))
	_ = e.WriteToken(jsontext.String("Title"))
	_ = e.WriteToken(jsontext.String(p.Title))
	_ = e.WriteToken(jsontext.String("Text"))
	_ = e.WriteToken(jsontext.String(p.Text))
	_ = e.WriteToken(jsontext.ObjectEnd)
	return nil
}

func (s *Stream) streamPagesWithMarshaler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	e := jsontext.NewEncoder(w)
	err = e.WriteToken(jsontext.ArrayStart)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}

	opts := jsonv2.WithMarshalers(jsonv2.MarshalFuncV2(marshalPage))
	for p, err := range s.db.StreamPages(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		err = jsonv2.MarshalEncode(e, &p, opts)
		if err != nil {
			s.logger.Error("fail to encode JSON", "err", err)
			return
		}
	}

	err = e.WriteToken(jsontext.ArrayEnd)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}
