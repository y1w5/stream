package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

func (s *Stream) listPages(w http.ResponseWriter, r *http.Request) {
	pages, err := s.service.ListPages(r.Context())
	if err != nil {
		s.logAndWriteError(w, "fail to list pages", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(pages)
	if err != nil {
		s.logger.Error("fail to encode JSON", "err", err)
	}
}

func (s *Stream) listPagesV2(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotImplemented
	http.Error(w, http.StatusText(status), status)

}

func (s *Stream) streamPagesV2(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotImplemented
	http.Error(w, http.StatusText(status), status)

}

func (s *Stream) logAndWriteError(w http.ResponseWriter, msg string, args ...any) {
	s.logger.Error(msg, args...)

	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	http.Error(w, http.StatusText(status), status)
}
