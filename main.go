package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/middlewares"
	"github.com/mohieey/lenslocked/migrations"
	"github.com/mohieey/lenslocked/models"
)

var port = ":3000"

func main() {
	// Setup the database
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

	err = models.MigrateFS(db, ".", migrations.FS)
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := models.UserService{DB: db}
	sessionService := models.SessionService{DB: db, BytesPerToken: 32}

	// Setup middlewares
	umw := middlewares.UserMiddleware{
		SessionService: &sessionService,
	}

	// csrfKey := "fivhenqpamrhdfgxtymopqfhmzxcarnd"
	// csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup controllers
	usersController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	// Setup routes
	r := chi.NewRouter()
	r.Use(umw.SetUser)

	r.Post("/signup", usersController.SignUp)
	r.Post("/signin", usersController.SignIn)
	r.Delete("/signout", usersController.SignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	//Start the server

	fmt.Println("Serving on ", port)
	http.ListenAndServe(port, r)
}
