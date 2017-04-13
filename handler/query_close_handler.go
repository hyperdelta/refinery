package handler

import (
	"net/http"
	"log"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/hyperdelta/refinery/config"
)

type QueryCloseHandler struct {
	Handler
}

/**
	constructor
 */
func NewQueryCloseHandler(r *mux.Router) *QueryCloseHandler {
	qch := new(QueryCloseHandler)
	qch.router = r
	return qch
}

/**
	register handler path
 */
func (h *QueryCloseHandler) RegisterHandlePath() {

	h.router.Handle("/_close/{id}", h)

	if config.Debug {
		log.Print("register path /_close/{id}, from QueryCloseHandler")
	}
}

func (h *QueryCloseHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	res, _, _, _ := jsonparser.Get(body,"test")
	rw.Write(res)
}

