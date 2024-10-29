package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest request.UserRequest
		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			panic(err)
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		user := userRequest.ToUser(string(password))
		err = service.CreateUser(user, db)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Fprint(w)
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds request.CredentialsRequest
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			panic(err)
		}

		err = service.VerifyCredentials(creds.Username, creds.Password, db)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := service.GenerateToken(creds.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		json.NewEncoder(w).Encode(resource.UserResource{Username: creds.Username})
	}
}

func AuthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(handler.UserIdKey).(int)
		fmt.Printf("userId: %v.\n", userId)
		user, err := service.FindUser(userId, db)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(resource.FromUser(user))
	}
}
