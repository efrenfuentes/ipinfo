package api

import (
	"log"
)

const version = "1.0.0"

type Config struct {
	Port int
	IPInfoKey  string
}

type API struct {
	Config Config
	Logger *log.Logger
}
