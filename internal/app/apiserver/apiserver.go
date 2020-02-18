package apiserver

import (
	"database/sql"
	//	"fmt"
	"net/http"
	//	"os"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
)

// Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	//	defer db.Close()
	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	srv := newServer(store, sessionStore)

	//	return http.ListenAndServe(config.BindAddr, srv)
	// http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images"))))
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	//func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	//	db, err := sql.Open("pgx", databaseURL)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
	//		os.Exit(1)
	//	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
