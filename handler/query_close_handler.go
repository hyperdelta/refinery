package handler

import (
	"net/http"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"io/ioutil"
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
	logger.Debug("register path /_close/{id}")
}

func (h *QueryCloseHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	res, _, _, _ := jsonparser.Get(body,"test")
	rw.Write(res)
}

