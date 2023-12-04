package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/views"
)

func excuteTemplate(w http.ResponseWriter, filePath string) {
	t, err := views.Parse(filePath)
	if err != nil {
		log.Printf("Error parsing %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	excuteTemplate(w, "home")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	excuteTemplate(w, "contact")

}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	excuteTemplate(w, "faq")

}

var port = ":3000"

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
