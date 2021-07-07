package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	infoLog *log.Logger
	errLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP socket address env variable")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		infoLog: infoLog,
		errLog: errLog,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))


	infoLog.Printf("Starting server on %s", *addr)

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errLog,
		Handler: mux,
	}
	err := server.ListenAndServe()

	errLog.Fatal(err)
}