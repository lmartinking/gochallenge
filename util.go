package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func ReadJsonBody(r *http.Request, v interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		return nil, err
	}
	if err = r.Body.Close(); err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
