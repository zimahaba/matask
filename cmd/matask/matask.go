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

	mux.HandleFunc("POST /signup", insecured(handler.MataskHandler{DB: db, F: transport.SignupHandler}))
	mux.HandleFunc("POST /auth/login", insecured(handler.MataskHandler{DB: db, F: transport.LoginHandler}))
	mux.HandleFunc("POST /auth/logout", transport.LogoutHandler)
	mux.HandleFunc("GET /auth/userinfo", secured(handler.MataskHandler{DB: db, F: transport.AuthCheckHandler}))

	mux.HandleFunc("GET /projects", secured(handler.MataskHandler{DB: db, F: transport.GetFilteredProjectsHandler}))
	mux.HandleFunc("GET /projects/{id}", secured(handler.MataskHandler{DB: db, F: transport.GetProjectHandler}))
	mux.HandleFunc("POST /projects", secured(handler.MataskHandler{DB: db, F: transport.CreateProjectHandler}))
	mux.HandleFunc("PUT /projects/{id}", secured(handler.MataskHandler{DB: db, F: transport.UpdateProjectHandler}))
	mux.HandleFunc("DELETE /projects/{id}", secured(handler.MataskHandler{DB: db, F: transport.DeleteProjectHandler}))

	mux.HandleFunc("GET /books/{id}", secured(handler.MataskHandler{DB: db, F: transport.GetBookHandler}))
	mux.HandleFunc("GET /books", secured(handler.MataskHandler{DB: db, F: transport.GetFilteredBooksHandler}))
	mux.HandleFunc("GET /books/cover/{id}", secured(handler.MataskHandler{DB: db, F: transport.GetBookCoverHandler}))
	mux.HandleFunc("POST /books", secured(handler.MataskHandler{DB: db, F: transport.SaveBookHandler}))
	mux.HandleFunc("PUT /books/{id}", secured(handler.MataskHandler{DB: db, F: transport.SaveBookHandler}))
	mux.HandleFunc("DELETE /books/{id}", secured(handler.MataskHandler{DB: db, F: transport.DeleteBookHandler}))

	mux.HandleFunc("GET /movies/{id}", secured(handler.MataskHandler{DB: db, F: transport.GetMovieHandler}))
	mux.HandleFunc("GET /movies", secured(handler.MataskHandler{DB: db, F: transport.GetFilteredMoviesHandler}))
	mux.HandleFunc("GET /movies/poster/{id}", secured(handler.MataskHandler{DB: db, F: transport.GetMoviePosterHandler}))
	mux.HandleFunc("POST /movies", secured(handler.MataskHandler{DB: db, F: transport.SaveMovieHandler}))
	mux.HandleFunc("PUT /movies/{id}", secured(handler.MataskHandler{DB: db, F: transport.SaveMovieHandler}))
	mux.HandleFunc("DELETE /movies/{id}", secured(handler.MataskHandler{DB: db, F: transport.DeleteMovieHandler}))

	mux.Handle("GET /tasks", secured(handler.MataskHandler{DB: db, F: transport.GetTasksHandler}))
	//mux.Handle("GET /tasks", handler.ErrorHandler{DB: db, H: h))

	mux.HandleFunc("PUT /images/{id}", secured(handler.MataskHandler{DB: db, F: transport.UploadImageHandler}))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
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

func secured(h handler.MataskHandler) http.HandlerFunc {
	return handler.Auth(h)
}

func insecured(h handler.MataskHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
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
