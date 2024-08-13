package main

import (
	"database/sql"
	"dkr290/go-advanced-projects/snippets-toolbox/internal/models"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var cfg config

	//define a new command-line flag with the name addr and default value
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:password@tcp(localhost)/snippetbox?parseTime=true", "MySql data source")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any
	// errors are
	// encountered during parsing the application will be terminated.

	flag.Parse()
	// use log.New for to create logger for writing informational message
	//prefix INFO or ERROR to stdout or stderr
	//additional info is local date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, using stderr
	// the destination and use the log.Lshortfile flag to include the
	// file name and line number.

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//database driver to keep main short openDB() function below
	db, err := openDB(cfg.dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	//close the connectionpool
	defer db.Close()

	//use template cache to be used
	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	// initialize new instance of application containing the dependencies

	app := &appconfig{
		errotLog:      errLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetsModel{DB: db},
		templateCache: templateCache,
	}
	// The value returned from the flag.String() function is a pointer to theflag
	// value, not the value itself. So we need to dereference the pointer
	// prefix it with the * symbol) before using it. Note that it is using
	// log.Printf() function to interpolate the address with the log message.

	//initialize a new http.Server sruct. We set the address and handler fields
	//the errorlog so the server is using a custom errolog logger

	srv := http.Server{
		Addr:     cfg.addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	if err := srv.ListenAndServe(); err != nil {

		errLog.Fatal(err)

	}
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	count := 0
	retryInterval := 2 * time.Second
	for {

		if err := db.Ping(); err == nil {
			log.Println("Sucesfully connected to the database")
			return db, nil
		} else {
			log.Printf("Attempt %d: Failed to connect to the database. Retrying in %v...\n", count, retryInterval)
			time.Sleep(retryInterval)
			count++
			if count > 10 {
				return nil, err
			}
		}
	}

}
