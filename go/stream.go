package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	_ "github.com/mattn/go-sqlite3"

	"github.com/y1w5/stream/go/internal/middleware"
)

// Stream is the main application. It stores and links all the components.
type Stream struct {
	server *http.Server
	db     *DB
	logger *slog.Logger

	errChan chan error
}

// NewStreamParams stores required parameters for [NewStream].
type NewStreamParams struct {
	Bind   string
	DB     string
	Logger *slog.Logger
}

// NewStream instanciates a [Stream].
func NewStream(arg NewStreamParams) (*Stream, error) {
	db, err := NewDB(arg.DB)
	if err != nil {
		return nil, fmt.Errorf("new db: %v", err)
	}

	return &Stream{
		db:      db,
		server:  &http.Server{Addr: arg.Bind},
		logger:  arg.Logger,
		errChan: make(chan error, 1),
	}, nil
}

// ListenAndServe listens and serves HTTP requests.
func (s *Stream) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/pages.list", s.listPages)
	mux.HandleFunc("/pages.stream", s.streamPages)
	s.server.Handler = middleware.Logger(s.logger, mux)

	go func() {
		s.logger.Info("listening on " + s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errChan <- fmt.Errorf("listen: %v", err)
			return
		}
		s.errChan <- nil
	}()

	select {
	case err := <-s.errChan:
		return err
	case <-time.After(time.Second):
		return nil
	}
}

// Close closes allocated ressources. It waits for 5 seconds for the running
// HTTP requests to finish before stopping.
func (s *Stream) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info("closing server & db")

	var errs []error
	if err := s.server.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("server: %v", err))
	}
	if err := <-s.errChan; err != nil {
		errs = append(errs, fmt.Errorf("server: %v", err))
	}
	if err := s.db.Close(); err != nil {
		errs = append(errs, fmt.Errorf("db: %v", err))
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

type response struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	http.Error(w, http.StatusText(status), status)
}

func (s *Stream) listPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	var pages []Page
	for p, err := range s.db.StreamPages(r.Context(), limit) {
		if err != nil {
			s.logger.Error("fail to stream pages", "err", err)
			return
		}
		pages = append(pages, p)
	}

	e := jsontext.NewEncoder(w)
	err = jsonv2.MarshalEncode(e, pages)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
		return
	}
}

func (s *Stream) streamPages(w http.ResponseWriter, r *http.Request) {
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

func parseLimit(r *http.Request) (int, error) {
	tmp := r.URL.Query().Get("limit")
	if tmp == "" {
		return 0, nil
	}
	limit, err := strconv.Atoi(tmp)
	if err != nil {
		return 0, nil
	}
	return limit, nil
}
