package main // define entry point

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"malaysia-lottery-checker/scraper"
)

type PageData struct {
	Title   string
	Message string
}

// entry function
func main() {
	indexTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to load HTML template: %v", err)
	}

	// to get static file path
	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.HandleFunc("/test/supreme-toto", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		result, err := scraper.FetchSupremeTotoResult()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("failed to encode Supreme Toto result: %v", err)
		}
	})

	//like controller route
	//w is used to send a response back to the browser.	r contains information about the incoming request.
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
