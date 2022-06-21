package main

import (
	"flag"
	"fmt"
	"github.com/efrenfuentes/ipinfo/internal/api"
	"github.com/efrenfuentes/ipinfo/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 4000, "API server port")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Setup all routes
	router := httprouter.New()
	router.GET("/", api.GetIpAddress)

	// Setup 404 / 405 error handlers
	router.NotFound = http.HandlerFunc(api.NotFound)
	router.MethodNotAllowed = http.HandlerFunc(api.MethodNotAllowed)

	routes := middleware.EnableCORS(router)

	// Start the server
	logger.Printf("starting API server on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), routes)
	logger.Fatal(err)
}
