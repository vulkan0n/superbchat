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
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/ui"
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

	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "frontend/dist", // This is the path to your SPA build folder, the folder that is created from running "npm build"
		Index:      "index.html",    // This is the default html page for your SPA
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(ui.Frontend),
	}))

	if err := e.Start(":8900"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}
	app.infoLog.Printf("App running in port: 8900")
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
