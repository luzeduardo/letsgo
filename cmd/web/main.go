package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"poc.eduardo-luz.eu/internal/models"
)

type config struct {
	addr      string
	staticDir string
}

// holds the app-wide dependencies for the application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	//making the models available to the handlers
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder //holds a pointer to the form.Decoder instance
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) // initi the pool for future use
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { //verify if is setup correctly creating a connection
		return nil, err
	}
	return db, nil
}

func main() {
	var cfg config
	dsn := flag.String("dsn", "", "MySQl data source name")

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	//log.New(destination, prefix, additional info)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	//initialize tempalte cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := *form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache, //adds to app deps
		formDecoder:   &formDecoder,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	// by default http server logs error to stdout
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg),
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
