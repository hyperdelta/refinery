package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/hyperdelta/refinery/pipeline"
)

type QueryCloseHandler struct {
	Handler
}

type QueryCloseResponseModel struct {
	Result string 		`json:"result"`
	Message string		`json:"message"`
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
	var params = mux.Vars(r)
	var id = params["id"]
	var response = QueryCloseResponseModel{};

	if id == "" {
		response.Result = "fail"
		response.Message = "ID need!"
	} else {
		pipeline.Close(id)
		// TODO: rethinkdb
		response.Result = "success"
	}

	ret, _ := json.Marshal(response)
	rw.Write(ret)
}

