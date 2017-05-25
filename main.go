package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperdelta/refinery/handler"
	"strconv"
	"github.com/hyperdelta/refinery/config"
	"github.com/hyperdelta/refinery/log"
)

var (
	conf = config.RefineryConfig
	defaultRouter *mux.Router
	logger *log.Logger = log.Get()
)

func main() {
	defaultRouter = mux.NewRouter()

	m := http.NewServeMux()
	handler.CreateDefaultRegisteredHandlers(defaultRouter)

	m.Handle("/", defaultRouter)
	http.DefaultServeMux = m

	logger.Info("Listen on " + conf.ListenAddress + strconv.Itoa(conf.ListenPort))

	http.ListenAndServe(conf.ListenAddress + ":" + strconv.Itoa(conf.ListenPort), nil)

}