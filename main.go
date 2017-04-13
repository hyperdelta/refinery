package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperdelta/refinery/handler"
	"strconv"
	"github.com/hyperdelta/refinery/config"
	"log"
)

var (
	conf = config.Config{ListenAddress: "", ListenPort:3000}
	defaultRouter *mux.Router
)

func main() {
	defaultRouter = mux.NewRouter()

	m := http.NewServeMux()
	handler.CreateDefaultRegisteredHandlers(defaultRouter)

	m.Handle("/", defaultRouter)
	http.DefaultServeMux = m


	if config.Debug {
		log.Print("Listen on " + conf.ListenAddress + strconv.Itoa(conf.ListenPort))
	}

	http.ListenAndServe(conf.ListenAddress + ":" + strconv.Itoa(conf.ListenPort), nil)

}