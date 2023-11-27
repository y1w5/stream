//go:build gotip

package main

import (
	"encoding/json"
	"net/http"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

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

	for p, err := range s.db.StreamPages(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		err = jsonv2.MarshalEncode(e, p)
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
