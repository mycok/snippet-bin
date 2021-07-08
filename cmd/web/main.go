package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mycok/snippet-bin/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog *log.Logger
	snippets *mysql.SnippetModel
}

func openDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP socket address env variable")
	dsn := flag.String("dsn", "webu:webu@/snippet_box?parseTime=true", "MysQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDBConnection(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	// only necessary if the application a graceful shutdown mechanism
	defer db.Close()

	app := &application{
		infoLog: infoLog,
		errLog: errLog,
		snippets: &mysql.SnippetModel{ DB: db},
	}

	infoLog.Printf("Starting server on %s", *addr)

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errLog,
		Handler: app.routes(),
	}
	err = server.ListenAndServe()

	errLog.Fatal(err)
}