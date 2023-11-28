//go:build !gotip

package main

import (
	"encoding/json"
	"net/http"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (s *Stream) listPagesV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	var pages []Page
	s.db.StreamPages(r.Context(), limit)(func(p Page, err error) bool {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return false
		}
		pages = append(pages, p)
		return true
	})

	e := jsontext.NewEncoder(w)
	err = jsonv2.MarshalEncode(e, pages)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

func (s *Stream) streamPagesV2(w http.ResponseWriter, r *http.Request) {
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

	s.db.StreamPages(r.Context(), limit)(func(p Page, err error) bool {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return false
		}
		err = jsonv2.MarshalEncode(e, p)
		if err != nil {
			s.logger.Error("fail to encode JSON", "err", err)
			return false
		}
		return true
	})

	err = e.WriteToken(jsontext.ArrayEnd)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

func (s *Stream) listPagesV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	var pages []Page
	s.db.StreamPageSlice(r.Context(), limit)(func(tmps []Page, err error) bool {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return false
		}
		pages = append(pages, tmps...)
		return true
	})

	e := jsontext.NewEncoder(w)
	err = jsonv2.MarshalEncode(e, pages)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

func (s *Stream) streamPagesV3(w http.ResponseWriter, r *http.Request) {
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

	s.db.StreamPageSlice(r.Context(), limit)(func(pages []Page, err error) bool {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return false
		}
		for _, p := range pages {
			err = jsonv2.MarshalEncode(e, p)
			if err != nil {
				s.logger.Error("fail to encode JSON", "err", err)
				return false
			}
		}
		return true
	})

	err = e.WriteToken(jsontext.ArrayEnd)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}
