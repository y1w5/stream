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
	mux.HandleFunc("/v1/pages.list.std", s.listPagesStd)
	mux.HandleFunc("/v1/pages.list.exp", s.listPagesExp)
	mux.HandleFunc("/v2/pages.list", s.listPagesV2)
	mux.HandleFunc("/v2/pages.stream", s.streamPagesV2)
	mux.HandleFunc("/v3/pages.list", s.listPagesV3)
	mux.HandleFunc("/v3/pages.stream", s.streamPagesV3)
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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	http.Error(w, http.StatusText(status), status)
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
