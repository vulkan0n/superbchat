package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"github.com/vulkan0n/superbchat/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	accounts       *models.AccountModel
	superchats     *models.SuperchatModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	dsn := flag.String("dsn", "postgres://web:pass@localhost/superbchat?sslmode=disable", "PostgreSQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := *form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		superchats:     &models.SuperchatModel{DB: db},
		accounts:       &models.AccountModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    &formDecoder,
		sessionManager: sessionManager,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8900"
	}

	srv := &http.Server{
		Addr:     ":" + port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("App running in port: %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
