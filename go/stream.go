package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	jsonv2 "github.com/go-json-experiment/json"
	_ "github.com/mattn/go-sqlite3"

	"github.com/y1w5/stream/go/internal/middleware"
)

// Stream is the main application. It stores and links all the components.
type Stream struct {
	server  *http.Server
	service *Service
	logger  *slog.Logger

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
	service, err := NewService(arg.DB)
	if err != nil {
		return nil, fmt.Errorf("new service: %v", err)
	}

	return &Stream{
		service: service,
		server:  &http.Server{Addr: arg.Bind},
		logger:  arg.Logger,
		errChan: make(chan error, 1),
	}, nil
}

// ListenAndServe listens and serves HTTP requests.
func (s *Stream) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/v1/pages.list", s.listPages)
	mux.HandleFunc("/v2/pages.list", s.listPagesV2)
	mux.HandleFunc("/v2/pages.stream", s.streamPagesV2)
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
	if err := s.service.Close(); err != nil {
		errs = append(errs, fmt.Errorf("service: %v", err))
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

func (s *Stream) listPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, err := s.service.ListPages(r.Context())
	if err != nil {
		s.logger.Error("fail to execute handler", "err", err)
		goto encode_err
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response{OK: true, Payload: u})
	return

encode_err:
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
}

func (s *Stream) listPagesV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, err := s.service.ListPages(r.Context())
	if err != nil {
		s.logger.Error("fail to execute handler", "err", err)
		goto encode_err
	}

	w.WriteHeader(http.StatusOK)
	_ = jsonv2.MarshalWrite(w, response{OK: true, Payload: u})
	return

encode_err:
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(response{Error: err.Error()})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	http.Error(w, http.StatusText(status), status)
}
