package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title   string
	Message string
}

func main() {
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to load HTML template: %v", err)
	}

	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := PageData{
			Title:   "Malaysia Lottery Checker",
			Message: "Go application is running inside Docker.",
		}

		if err := indexTemplate.Execute(w, data); err != nil {
			log.Printf("failed to render page: %v", err)
			http.Error(w, "failed to render page", http.StatusInternalServerError)
		}
	})

	address := ":8080"

	log.Printf("server running at http://localhost%s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}