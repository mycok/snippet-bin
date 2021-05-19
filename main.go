package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}
	w.Write([]byte("Hello from snippet box...."))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet...."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed!", 405)
	
		return
	}

	w.Write([]byte("Create a new snippet...."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}