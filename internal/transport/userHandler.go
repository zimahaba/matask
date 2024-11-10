package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userRequest request.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := userRequest.ToUser(string(password))
	err = service.CreateUser(user, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var creds request.CredentialsRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.VerifyCredentials(creds.Username, creds.Password, db)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenCookie, err := service.GenerateTokenCookie(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, tokenCookie)

	if creds.KeepLoggedIn {
		refreshCookie, err := service.GenerateRefreshCookie(creds.Username, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, refreshCookie)
	}

	json.NewEncoder(w).Encode(resource.UserResource{Username: creds.Username})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, service.GenerateCookie(service.TOKEN_COOKIE_NAME, "", -1))
	http.SetCookie(w, service.GenerateCookie(service.REFRESH_COOKIE_NAME, "", -1))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func RefreshHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	refreshCookie, err := r.Cookie(service.REFRESH_COOKIE_NAME)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := service.FindUsernameByRefreshToken(refreshCookie.Value, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenCookie, err := service.GenerateTokenCookie(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, tokenCookie)

	refreshCookie, err = service.GenerateRefreshCookie(username, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, refreshCookie)

	json.NewEncoder(w).Encode(resource.UserResource{Username: username})
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, redis *redis.Client) {
	var creds request.CredentialsRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user exists
	userId, err := service.FindUserId(creds.Username, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	randomToken, err := service.GenerateRandomToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("generated token: %v", randomToken)

	err = redis.Set(context.Background(), randomToken, userId, 24*time.Second).Err()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send email
}

func RecoverPasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, redis *redis.Client) {
	token := r.URL.Query()["tk"][0]
	if token == "" {
		http.Error(w, "Token is required.", http.StatusBadRequest)
		return
	}

	userId, err := redis.Get(context.Background(), token).Result()

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("userId %v", userId)

}

func UserInfoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId := r.Context().Value(handler.UserIdKey).(int)
	user, err := service.FindUser(userId, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(resource.FromUser(user))
}
