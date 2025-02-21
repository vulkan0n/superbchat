package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/vulkan0n/superbchat/internal/models"
)

type application struct {
	echo       *echo.Echo
	infoLog    *log.Logger
	errorLog   *log.Logger
	accounts   *models.AccountModel
	superchats *models.SuperchatModel
	wsHub      *WebSocketHub
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

	wsHub := newWebSocketHub()
	e := echo.New()

	app := &application{
		echo:       e,
		errorLog:   errorLog,
		infoLog:    infoLog,
		superchats: &models.SuperchatModel{DB: db},
		accounts:   &models.AccountModel{DB: db},
		wsHub:      wsHub,
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
