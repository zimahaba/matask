package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"matask/internal/transport"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	loadEnv()
	db = connectDB()
	defer db.Close()

	// https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/
	// 	https://www.willem.dev/articles/url-path-parameters-in-routes/
	mux := http.NewServeMux()
	mux.HandleFunc("GET /projects/{id}", transport.GetProjectHandler(db))
	mux.HandleFunc("POST /projects", transport.CreateProjectHandler(db))
	mux.HandleFunc("PUT /projects/{id}", transport.UpdateProjectHandler(db))

	mux.HandleFunc("GET /books/{id}", transport.GetBookHandler(db))
	mux.HandleFunc("POST /books", transport.CreateBookHandler(db))
	mux.HandleFunc("PUT /books/{id}", transport.UpdateBookHandler(db))

	mux.HandleFunc("GET /movies/{id}", transport.GetMovieHandler(db))
	mux.HandleFunc("POST /movies", transport.CreateMovieHandler(db))
	mux.HandleFunc("PUT /movies/{id}", transport.UpdateMovieHandler(db))

	mux.HandleFunc("GET /tasks", transport.GetTasksHandler(db))

	mux.HandleFunc("POST /cover", transport.UploadCoverHandler())

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
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
