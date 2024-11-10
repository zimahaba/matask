package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const (
	TOKEN_COOKIE_NAME   = "token"
	REFRESH_COOKIE_NAME = "refresh"
	tokenExpiration     = 60 * time.Minute
	cookieExpiration    = 86400 // seconds
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateTokenCookie(username string) (*http.Cookie, error) {
	token, err := generateToken(username)
	if err != nil {
		slog.Error(err.Error())
		return &http.Cookie{}, err
	}

	return GenerateCookie(TOKEN_COOKIE_NAME, token, cookieExpiration), nil
}

func GenerateRefreshCookie(username string, db *sql.DB) (*http.Cookie, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		slog.Error(err.Error())
		return &http.Cookie{}, err
	}

	err = UpsertRefreshToken(refreshToken, username, db)
	if err != nil {
		slog.Error(err.Error())
		return &http.Cookie{}, err
	}

	return GenerateCookie(REFRESH_COOKIE_NAME, refreshToken, math.MaxInt32), nil
}

func GenerateCookie(name string, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteStrictMode,
	}
}

func generateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(JwtKey)
}

func generateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
