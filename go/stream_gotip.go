//go:build gotip

package main

import (
	"net/http"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (s *Stream) streamPagesV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	e := jsontext.NewEncoder(w)
	err := e.WriteToken(jsontext.ArrayStart)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
	for p, err := range s.service.StreamPages(r.Context()) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			break
		}
		err = jsonv2.MarshalEncode(e, p)
		if err != nil {
			s.logger.Error("fail to encode JSON", "err", err)
			break
		}
	}
	err = e.WriteToken(jsontext.ArrayEnd)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}
