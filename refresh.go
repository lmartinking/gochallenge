package main

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Success     bool   `json:"success"`
	Error       error  `json:"error"`
	AccessToken string `json:"access_token"`
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var refreshReq RefreshRequest
	if _, err := ReadJsonBody(r, &refreshReq); err != nil {
		return
	}

	claims, err := VerifyToken(refreshReq.RefreshToken, Cfg.PubKey)
	if err != nil {
		log.Error("Error verifying token:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Debug("Claims:", claims)

	w.Header().Set("Content-Type", "application/json")

	var refreshResp RefreshResponse
	refreshResp.Success = false

	if claims.TokenType != "refresh" {
		refreshResp.Error = errors.New("Invalid token type")
	} else {
		acct := GetAccount(claims.User)
		tok, err := MakeAccessToken(*acct, Cfg.PrivKey)
		if err != nil {
			refreshResp.Error = errors.New("Could not generate new token")
		} else {
			refreshResp.AccessToken = tok
			refreshResp.Success = true
		}
	}

	if refreshResp.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(refreshResp)
}
