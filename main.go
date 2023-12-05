package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/templates"
	"github.com/mohieey/lenslocked/views"
)

var port = ":3000"

func main() {
	r := chi.NewRouter()

	tmpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout_parts.gohtml"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "layout_parts.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "layout_parts.gohtml"))
	r.Get("/faq", controllers.FAQ(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
