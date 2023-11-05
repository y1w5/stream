// This package implements a streaming pipeline between a database and a controller.
package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
)

// User represents an user in our awesome banking app.
type User struct {
	Name    string
	Balance int
}

// db mocks our database.
var db = []User{
	{Name: "Gallant Bassi", Balance: 40282},
	{Name: "Stupefied Lederberg", Balance: 34934},
	{Name: "Eloquent Proskuriakova", Balance: 42924},
	{Name: "Admiring Chebyshev", Balance: 37802},
	{Name: "Funny Jang", Balance: 19530},
	{Name: "Adoring Mayer", Balance: 38438},
	{Name: "Jolly Noether", Balance: 17229},
	{Name: "Wizardly Fermat", Balance: 48600},
	{Name: "Agitated Einstein", Balance: 14764},
	{Name: "Stoic Black", Balance: 25202},
	{Name: "Mystifying Blackwell", Balance: 49565},
	{Name: "Youthful Maxwell", Balance: 16125},
	{Name: "Musing Sutherland", Balance: 22296},
	{Name: "Friendly Almeida", Balance: 14602},
	{Name: "Vigilant Cray", Balance: 37680},
	{Name: "Exciting Zhukovsky", Balance: 14975},
	{Name: "Naughty Keldysh", Balance: 20369},
	{Name: "Practical Noether", Balance: 6804},
	{Name: "Zen Lewin", Balance: 23631},
	{Name: "Clever Vaughan", Balance: 26239},
}

// StreamUsersFromDB streams users from the database.
//
// TODO: use slice of User for improved performance.
func StreamUsersFromDB(fn func(*User) error) error {
	for _, u := range db {
		if err := fn(&u); err != nil {
			return err
		}
	}
	return nil
}

// StreamUsersFromService streams users from the database, with added business logic.
func StreamUsersFromService(fn func(*User) error) error {
	return StreamUsersFromDB(func(u *User) error {
		// We remove the poor, they are not paying enough for our services.
		if u.Balance < 20_000 {
			return nil
		}

		return fn(u)
	})
}

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
