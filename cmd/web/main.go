package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mycok/snippet-bin/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	session *sessions.Session
	infoLog *log.Logger
	errLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
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
	secret := flag.String("secret", "yeueuu+hffs24453+42fggsg*yu@etyr", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDBConnection(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	// only necessary if the application a graceful shutdown mechanism
	defer db.Close()
	// cache all template pages when the application starts
	tempCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		infoLog: infoLog,
		errLog: errLog,
		snippets: &mysql.SnippetModel{ DB: db},
		templateCache: tempCache,
		session: session,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	s := &http.Server{
		Addr: *addr,
		TLSConfig: tlsConfig,
		ErrorLog: errLog,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = s.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	errLog.Fatal(err)
}