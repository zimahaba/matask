package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"matask/internal/service"
	"matask/internal/transport/request"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))

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

		if creds.Username != "admin" || creds.Password != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(1 * time.Minute)
		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: time.Now().Add(5 * time.Minute),
		})
	}
}
