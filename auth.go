package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	username string
	password string
}

var accounts = []Account{}

func makePasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

func comparePasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func makeAccount(username string, password string) Account {
	acct := Account{username: username, password: makePasswordHash(password)}
	log.Debug("makeAccount:", acct)
	return acct
}

func GetAccount(username string) *Account {
	for _, acct := range accounts {
		if acct.username == username {
			return &acct
		}
	}
	return nil
}

func hasAccount(username string) bool {
	return GetAccount(username) != nil
}

func AddAccount(username string, password string) bool {
	if hasAccount(username) {
		return false
	}

	accounts = append(accounts, makeAccount(username, password))

	return true
}

func AuthAccount(username string, password string) bool {
	if username == "" {
		return false
	}
	acct := GetAccount(username)
	if acct == nil {
		return false
	}
	return comparePasswordHash(acct.password, password)
}
