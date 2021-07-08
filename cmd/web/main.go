package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog *log.Logger
	errLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP socket address env variable")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog: infoLog,
		errLog: errLog,
	}

	infoLog.Printf("Starting server on %s", *addr)

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errLog,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()

	errLog.Fatal(err)
}