package api

import (
	"fmt"
	"net/http"
)

func (api *API) LogError(r *http.Request, err error) {
	api.Logger.Printf("%s %s %s", r.Method, r.URL.Path, err)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func (api *API) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := Envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := api.WriteJSON(w, status, env, nil)
	if err != nil {
		api.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func (api *API) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.LogError(r, err)

	message := "the server encountered a problem and could not process your request"
	api.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send 404 Not Found status code and
// JSON response to the client.
func (api *API) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the request resource could not be found"
	api.ErrorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to client
func (api *API) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	api.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}
