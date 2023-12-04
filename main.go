package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/views"
)

var port = ":3000"

func main() {
	r := chi.NewRouter()

	tmpl := views.Must(views.Parse("home"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.Parse("contact"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.Parse("faq"))
	r.Get("/faq", controllers.StaticHandler(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
