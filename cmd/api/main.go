package main

import (
	"flag"
	"github.com/efrenfuentes/ipinfo/internal/api"
	"log"
	"os"
)

func main() {
	var cfg api.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.IpinfoAccessToken, "ipinfo_access_token", "", "IPInfo Access Token")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	api := &api.API{
		Config: cfg,
		Logger: logger,
	}

	err := api.Serve()
	if err != nil {
		logger.Fatal(err) 
	}
}
