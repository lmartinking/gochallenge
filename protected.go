package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func extractTokenFromAuthHeader(v string) (string, error) {
	if strings.HasPrefix(v, "bearer ") {
		parts := strings.Split(v, " ")
		if len(parts) == 2 {
			tok := strings.Replace(parts[1], " ", "", -1)
			return tok, nil
		}
	}
	return "", errors.New("invalid auth header value")
}

func ProtectedWrapper(handler func(http.ResponseWriter, *http.Request, *TokenClaims)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := extractTokenFromAuthHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, err := VerifyToken(tok, Cfg.PubKey)
		if err != nil {
			log.Error("Token verification failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r, claims)
	}
}

func HandleProtected(w http.ResponseWriter, r *http.Request, claims *TokenClaims) {
	fmt.Fprintf(w, "Hello %s", claims.User)
}
