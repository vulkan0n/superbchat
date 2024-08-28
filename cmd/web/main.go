package main

import (
	"database/sql"
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/vulkan0n/superbchat/internal/models"
)

type application struct {
	echo           *echo.Echo
	infoLog        *log.Logger
	errorLog       *log.Logger
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

	e := echo.New()

	app := &application{
		echo:           e,
		errorLog:       errorLog,
		infoLog:        infoLog,
		superchats:     &models.SuperchatModel{DB: db},
		accounts:       &models.AccountModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    &formDecoder,
		sessionManager: sessionManager,
	}
	app.routes()

	if err := app.echo.Start(":8900"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.errorLog.Println(err)
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
