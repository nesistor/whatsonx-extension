package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"log"
)

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return errors.New("unable to parse JSON")
	}
	return nil
}

// writeJSON writes JSON responses with proper content type
func (app *Config) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		return errors.New("unable to write JSON")
	}
	return nil
}

// errorJSON writes error messages in JSON format
func (app *Config) errorJSON(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)

	var jsonResp jsonResponse
	jsonResp.Error = true
	jsonResp.Message = err.Error()

	app.writeJSON(w, statusCode, jsonResp)
}
