package main

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success      bool   `json:"success"`
	Error        error  `json:"error"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if _, err := ReadJsonBody(r, &loginRequest); err != nil {
		return
	}

	log.Debug("Login request:", loginRequest)

	var loginResponse LoginResponse
	loginResponse.Success = false

	authOk := AuthAccount(loginRequest.Username, loginRequest.Password)

	if authOk {
		account := GetAccount(loginRequest.Username)
		tok, err := MakeAccessToken(*account, Cfg.PrivKey)
		if err != nil {
			loginResponse.Error = err
		} else {
			loginResponse.AccessToken = tok
		}
		tok, err = MakeRefreshToken(*account, Cfg.PrivKey)
		if err != nil {
			loginResponse.Error = err
		} else {
			loginResponse.RefreshToken = tok
			loginResponse.Success = true
		}
	} else {
		loginResponse.Error = errors.New("auth failed")
	}

	w.Header().Set("Content-Type", "application/json")

	if loginResponse.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(loginResponse)
}
