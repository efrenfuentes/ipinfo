package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/efrenfuentes/ipinfo/internal/middleware"
)

func (api *API) Routes() http.Handler {
	router := httprouter.New()

	// Convert the notFoundResponse() helper to http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(api.NotFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", api.GetIpHandler)
	router.HandlerFunc(http.MethodGet, "/healthcheck", api.HealthcheckHandler)

	return middleware.EnableCORS(router)
}