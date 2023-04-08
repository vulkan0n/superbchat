package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vulkan0n/superbchat/internal/models"
)

var BCHAddress = ""
var ScamThreshold = 0.0001 // MINIMUM DONATION AMOUNT
var MessageMaxChar = 250
var NameMaxChar = 25
var username = "admin" // chat log /view page
var password = "adminadmin"
var AlertWidgetRefreshInterval = "10" //seconds
var checked = ""

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	superchats    *models.SuperchatModel
	templateCache map[string]*template.Template
}

func main() {

	conf := getDefaultConfig()
	BCHAddress = conf.BCHAddress
	ScamThreshold = conf.MinimumDonation
	MessageMaxChar = conf.MaxMessageChars
	NameMaxChar = conf.MaxNameChars
	username = conf.WebViewUsername
	password = conf.WebViewPassword
	AlertWidgetRefreshInterval = conf.OBSWidgetRefresh
	if conf.Checked == true {
		checked = " checked"
	}

	flag.StringVar(&BCHAddress, "addr", "bitcoincash:address", "Bitcoin Cash address to recieve founds")
	dsn := flag.String("dsn", "web:pass@/superbchat?parseTime=true", "MySQL data source name")
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

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		superchats:    &models.SuperchatModel{DB: db},
		templateCache: templateCache,
	}

	infoLog.Println(BCHAddress)

	infoLog.Println(fmt.Sprintf("OBS Alert path: /alert?auth=%s", password))

	// Create files and directory if they don't exist
	logDirectory := "./cmd/log"
	_ = os.Mkdir(logDirectory, os.ModePerm)

	_, err = os.OpenFile(logDirectory+"/paid.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorLog.Fatal(err)
	}

	_, err = os.OpenFile(logDirectory+"/alertqueue.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorLog.Fatal(err)
	}

	_, err = os.OpenFile(logDirectory+"/superchats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorLog.Fatal(err)
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

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
