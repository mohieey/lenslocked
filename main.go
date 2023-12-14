package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/models"
	"github.com/mohieey/lenslocked/templates"
	"github.com/mohieey/lenslocked/views"
)

var port = ":3000"

func main() {
	r := chi.NewRouter()

	tmpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tmpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected to the database successfully")
	defer db.Close()

	userService := models.UserService{DB: db}

	usersController := controllers.Users{UserService: &userService}
	r.Post("/signup", usersController.SignUp)
	r.Post("/signin", usersController.SignIn)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
