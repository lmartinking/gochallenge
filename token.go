package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type TokenClaims struct {
	TokenType string `json:"type"`
	User      string `json:"user"`
	jwt.StandardClaims
}

func LoadPrivKey(path string) (*rsa.PrivateKey, error) {
	priv_key_data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rsa_priv_key, err := jwt.ParseRSAPrivateKeyFromPEM(priv_key_data)
	if err != nil {
		return nil, err
	}

	return rsa_priv_key, nil
}

func LoadPubKey(path string) (*rsa.PublicKey, error) {
	pub_key_data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rsa_pub_key, err := jwt.ParseRSAPublicKeyFromPEM(pub_key_data)
	if err != nil {
		return nil, err
	}

	return rsa_pub_key, nil
}

func AccessTokenClaimsForAccount(account Account) TokenClaims {
	return TokenClaims{
		"access",
		account.username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24,
			Issuer:    "gochallenge",
		},
	}
}

func RefreshTokenClaimsForAccount(account Account) TokenClaims {
	return TokenClaims{
		"refresh",
		account.username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24*7,
			Issuer:    "gochallenge",
		},
	}
}

func MakeAccessToken(account Account, key interface{}) (string, error) {
	claims := AccessTokenClaimsForAccount(account)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signed, err := token.SignedString(key)
	log.Debug("access token:", signed)
	return signed, err
}

func MakeRefreshToken(account Account, key interface{}) (string, error) {
	claims := RefreshTokenClaimsForAccount(account)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signed, err := token.SignedString(key)
	log.Debug("refresh token:", signed)
	return signed, err
}

func VerifyToken(token_string string, key interface{}) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(token_string, claims, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		log.Error("Error parsing:", err)
		return nil, err
	}

	if token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
