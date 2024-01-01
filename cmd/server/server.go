package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lpernett/godotenv"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/middlewares"
	"github.com/mohieey/lenslocked/migrations"
	"github.com/mohieey/lenslocked/models"
	"github.com/mohieey/lenslocked/services"
)

type config struct {
	PSQL   services.PostgresConfig
	SMTP   models.SMPTPConfig
	Server struct {
		Host string
		Port string
	}
}

func loadEnvConfig() (*config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading environment configuration file: %v", err)
	}

	var config config

	config.Server.Host = os.Getenv("HOST")
	config.Server.Port = os.Getenv("PORT")

	config.PSQL.Host = os.Getenv("PSQL_HOST")
	config.PSQL.Port = os.Getenv("PSQL_PORT")
	config.PSQL.User = os.Getenv("PSQL_USER")
	config.PSQL.Password = os.Getenv("PSQL_PASSWORD")
	config.PSQL.DbName = os.Getenv("PSQL_NAME")
	config.PSQL.SSLMode = os.Getenv("PSQL_SSLMODE")

	if config.PSQL.Host == "" {
		config.PSQL = services.DefaultPostgresConfig()
	}

	config.SMTP.Host = os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	config.SMTP.Port = smtpPort
	config.SMTP.Username = os.Getenv("SMTP_USERNAME")
	config.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	return &config, nil
}

func main() {
	// Setup the database
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := services.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected to the database successfully")
	defer db.Close()

	err = services.MigrateFS(db, ".", migrations.FS)
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := &services.UserService{
		DB: db,
	}
	sessionService := &services.SessionService{
		DB: db, BytesPerToken: 32,
	}
	pwResetService := &services.PasswordResetService{
		DB:            db,
		BytesPerToken: 32,
	}
	emailService := services.NewEmailService(cfg.SMTP)
	galleryService := services.GalleryService{
		DB:                    db,
		SupportedExtensions:   []string{".jpeg", ".jpg", ".png", ".gif"},
		SupportedContentTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}

	// Setup middlewares
	umw := middlewares.UserMiddleware{
		SessionService: sessionService,
	}

	// csrfKey := "fivhenqpamrhdfgxtymopqfhmzxcarnd"
	// csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup controllers
	usersController := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	galleriesController := controllers.Galleries{
		GalleryService: &galleryService,
	}

	// Setup routes
	r := chi.NewRouter()
	r.Use(umw.SetUser)

	r.Post("/signup", usersController.SignUp)
	r.Post("/signin", usersController.SignIn)
	r.Delete("/signout", usersController.SignOut)
	r.Post("/forgot_password", usersController.ForgotPassword)
	r.Post("/reset_password", usersController.ResetPassord)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesController.Show)
		r.Get("/{id}/images/{filename}", galleriesController.ServeImage)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesController.Index)
			r.Post("/", galleriesController.Create)
			r.Put("/{id}", galleriesController.Update)
			r.Delete("/{id}", galleriesController.Delete)
			r.Post("/{id}/images", galleriesController.UploadImage)
			r.Delete("/{id}/images/{filename}", galleriesController.DeleteImage)
		})

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	//Start the server
	address := fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("Serving on ", address)
	err = http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
