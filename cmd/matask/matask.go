package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"matask/internal/transport"
	"matask/internal/transport/handler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/
// 	https://www.willem.dev/articles/url-path-parameters-in-routes/

// https://blog.questionable.services/article/http-handler-error-handling-revisited/
// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
// https://dev.to/neelp03/adding-logging-and-error-handling-middleware-to-your-go-api-2f33
// https://medium.com/geekculture/learn-go-middlewares-by-examples-da5dc4a3b9aa
// https://drstearns.github.io/tutorials/gomiddleware/
// https://www.jetbrains.com/guide/go/tutorials/authentication-for-go-apps/auth/

var db *sql.DB

func main() {
	loadEnv()
	db = connectDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", transport.SignupHandler(db))
	mux.HandleFunc("POST /login", transport.LoginHandler(db))

	mux.HandleFunc("GET /projects/{id}", auth(transport.GetProjectHandler(db)))
	mux.HandleFunc("POST /projects", auth(transport.CreateProjectHandler(db)))
	mux.HandleFunc("PUT /projects/{id}", auth(transport.UpdateProjectHandler(db)))
	mux.HandleFunc("DELETE /projects/{id}", auth(transport.DeleteProjectHandler(db)))

	mux.HandleFunc("GET /books/{id}", auth(transport.GetBookHandler(db)))
	mux.HandleFunc("POST /books", auth(transport.CreateBookHandler(db)))
	mux.HandleFunc("PUT /books/{id}", auth(transport.UpdateBookHandler(db)))
	mux.HandleFunc("DELETE /books/{id}", auth(transport.DeleteBookHandler(db)))

	mux.HandleFunc("GET /movies/{id}", auth(transport.GetMovieHandler(db)))
	mux.HandleFunc("POST /movies", auth(transport.CreateMovieHandler(db)))
	mux.HandleFunc("PUT /movies/{id}", auth(transport.UpdateMovieHandler(db)))
	mux.HandleFunc("DELETE /movies/{id}", auth(transport.DeleteMovieHandler(db)))

	mux.Handle("GET /tasks", auth(transport.GetTasksHandler(db)))
	//mux.Handle("GET /tasks", handler.ErrorHandler{DB: db, H: h))

	mux.HandleFunc("PUT /images/{id}", auth(transport.UploadImageHandler(db)))

	newMux := handler.Logging(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: newMux,
	}

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}

func auth(h http.Handler) http.HandlerFunc {
	return handler.Auth(h)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectDB() *sql.DB {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully connected to database")
	return db
}
