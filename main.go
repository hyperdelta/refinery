package main

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/TykTechnologies/goagain"
	"time"
	"os"
	"fmt"
)

var (
	config = Config{}
	NodeID string
	defaultRouter *mux.Router
)

func generateListener(listenAddress string, listenPort int) (net.Listener, error) {
	if listenAddress == "" {
		listenAddress = config.ListenAddress
	}

	if listenPort == 0 {
		listenPort = config.ListenPort
	}

	laddr := fmt.Sprintf("%s:%d", listenAddress, listenPort)

	return net.Listen("tcp", laddr)
}

func listen(l net.Listener, laddr string) {
	readTimeout := 120
	writeTimeout := 120

	if config.HttpServerOptions.ReadTimeout > 0 {
		readTimeout = config.HttpServerOptions.ReadTimeout
	}

	if config.HttpServerOptions.WriteTimeout > 0 {
		writeTimeout = config.HttpServerOptions.WriteTimeout
	}

	s := &http.Server{
		Addr: ":" + laddr,
		ReadTimeout: time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}

	newServeMux := http.NewServeMux()
	newServeMux.Handle("/", defaultRouter)
	http.DefaultServeMux = newServeMux

	go s.Serve(l)

}

func setDefaultRouter() {
	if defaultRouter == nil {
		defaultRouter = mux.NewRouter()
	}

	defaultRouter.Methods("GET", "POST")

}

var amForked bool

/**
Hook function before terminate
 */
func onFork() {
		//log.Info("Waiting to de-register")
		time.Sleep(10 * time.Second)

		os.Setenv("TYK_SERVICE_NODEID", NodeID)


	amForked = true
}

func main() {
	NodeID = generateRandomNodeID()
	l, goAgainErr := goagain.Listener(onFork)

	setDefaultRouter()

	if goAgainErr != nil {
		var err error
		if l, err = generateListener("", 0); err != nil {
			// error log
		}

		listen(l, l.Addr().String())

	}

}