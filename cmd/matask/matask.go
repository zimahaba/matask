package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"matask/internal/transport"
	"matask/internal/transport/handler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/
// 	https://www.willem.dev/articles/url-path-parameters-in-routes/

// https://blog.questionable.services/article/http-handler-error-handling-revisited/
// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
// https://dev.to/neelp03/adding-logging-and-error-handling-middleware-to-your-go-api-2f33
// https://medium.com/geekculture/learn-go-middlewares-by-examples-da5dc4a3b9aa
// https://drstearns.github.io/tutorials/gomiddleware/
// https://www.jetbrains.com/guide/go/tutorials/authentication-for-go-apps/auth/

// https://github.com/golang/example/blob/master/slog-handler-guide/README.md

func getLogLevel() slog.Level {
	level := os.Getenv("LOG_LEVEL")
	switch {
	case level == "DEBUG":
		return slog.LevelDebug
	case level == "INFO":
		return slog.LevelInfo
	case level == "WARN":
		return slog.LevelWarn
	case level == "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	loadEnv()
	db := connectDB()
	defer db.Close()

	slog.SetLogLoggerLevel(getLogLevel())

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", transport.SignupHandler(db))
	mux.HandleFunc("POST /login", transport.LoginHandler(db))
	mux.HandleFunc("POST /logout", transport.LogoutHandler())

	mux.HandleFunc("GET /auth/status", auth(transport.AuthCheckHandler(db), db))

	mux.HandleFunc("GET /projects/{id}", auth(transport.GetProjectHandler(db), db))
	mux.HandleFunc("POST /projects", auth(transport.CreateProjectHandler(db), db))
	mux.HandleFunc("PUT /projects/{id}", auth(transport.UpdateProjectHandler(db), db))
	mux.HandleFunc("DELETE /projects/{id}", auth(transport.DeleteProjectHandler(db), db))

	mux.HandleFunc("GET /books/{id}", auth(transport.GetBookHandler(db), db))
	mux.HandleFunc("GET /books", auth(transport.GetFilteredBooksHandler(db), db))
	mux.HandleFunc("GET /books/cover/{id}", auth(transport.GetBookCoverHandler(db), db))
	mux.HandleFunc("POST /books", auth(transport.CreateBookHandler(db), db))
	mux.HandleFunc("PUT /books/{id}", auth(transport.UpdateBookHandler(db), db))
	mux.HandleFunc("DELETE /books/{id}", auth(transport.DeleteBookHandler(db), db))

	mux.HandleFunc("GET /movies/{id}", auth(transport.GetMovieHandler(db), db))
	mux.HandleFunc("POST /movies", auth(transport.CreateMovieHandler(db), db))
	mux.HandleFunc("PUT /movies/{id}", auth(transport.UpdateMovieHandler(db), db))
	mux.HandleFunc("DELETE /movies/{id}", auth(transport.DeleteMovieHandler(db), db))

	mux.Handle("GET /tasks", auth(transport.GetTasksHandler(db), db))
	//mux.Handle("GET /tasks", handler.ErrorHandler{DB: db, H: h))

	mux.HandleFunc("PUT /images/{id}", auth(transport.UploadImageHandler(db), db))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-type"},
		AllowCredentials: true,
	})

	newMux := c.Handler(handler.Logging(mux))

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

func auth(h http.Handler, db *sql.DB) http.HandlerFunc {
	return handler.Auth(h, db)
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
