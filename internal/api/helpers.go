package api

import (
	"encoding/json"
	"net/http"
)

// Define an evelope type.
type Envelope map[string]interface{}

// Define WritePlainText helper for sending responses.
func (api *API) WritePlainText(w http.ResponseWriter, status int, body string, headers http.Header) error {
	// Write the headers.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: text/plain" header.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, err := w.Write([]byte(body))
	return err
}

// Define a WriteJSON helper for sending responses.
func (api *API) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	// Encode the data to JSON.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal.
	js = append(js, '\n')

	// Write the headers.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}