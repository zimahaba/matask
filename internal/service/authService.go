package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log/slog"
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
	cookieExpiration    = 10 * time.Second
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateTokenCookie(username string) (*http.Cookie, error) {
	token, err := generateToken(username)
	if err != nil {
		slog.Error(err.Error())
		return &http.Cookie{}, err
	}

	return GenerateCookie(TOKEN_COOKIE_NAME, token, time.Now().Add(cookieExpiration)), nil
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

	return GenerateCookie(REFRESH_COOKIE_NAME, refreshToken, time.Time{}), nil
}

func GenerateCookie(name string, value string, expires time.Time) *http.Cookie {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if !expires.IsZero() {
		c.Expires = expires
	}

	return c
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
