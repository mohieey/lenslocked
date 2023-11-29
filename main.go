package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	data, _ := json.Marshal("sss")
	w.Write(data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1>Contact Page</h1>
	<a href="mailto:mohieey@gmail.com">mohieey@gmail.com</a>`)
}

var port = ":3000"

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
