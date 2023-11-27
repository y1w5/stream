// This package implements a streaming pipeline between a database and a controller.
package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	params := NewStreamParams{
		Logger: slog.Default(),
	}
	flag.StringVar(&params.Bind, "bind", "127.0.0.1:8080", "adress of the HTTP server")
	flag.StringVar(&params.DB, "db", "stream.db", "path to the SQLite database")
	flag.Parse()

	stream, err := NewStream(params)
	if err != nil {
		fatal("fail to instanciate Stream", "err", err)
	}

	err = stream.ListenAndServe()
	if err != nil {
		fatal("fail to listen and serve", "err", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	err = stream.Close()
	if err != nil {
		fatal("fail to close Stream", "err", err)
	}
}

func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
