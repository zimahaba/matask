package transport

import (
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

	fmt.Printf("keep? %v.\n", creds.KeepLoggedIn)
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
	http.SetCookie(w, service.GenerateCookie(service.TOKEN_COOKIE_NAME, "", time.Unix(0, 1)))
	http.SetCookie(w, service.GenerateCookie(service.REFRESH_COOKIE_NAME, "", time.Unix(0, 1)))
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

func AuthCheckHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId := r.Context().Value(handler.UserIdKey).(int)
	user, err := service.FindUser(userId, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(resource.FromUser(user))
}
