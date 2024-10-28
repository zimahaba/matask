package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"matask/internal/service"
	"matask/internal/transport/request"
	"net/http"

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

		h := w.Header()
		h.Set("Authorization", token)
		w.WriteHeader(200)
	}
}
