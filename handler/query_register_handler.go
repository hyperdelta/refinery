package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"github.com/hyperdelta/refinery/query"
	"github.com/hyperdelta/refinery/pipeline"
)

type QueryRegisterHandler struct {
	Handler
}

type QueryRegisterResponseModel struct {
	Result string 		`json:"result"`
	Message string		`json:"message"`
	Endpoint string		`json:"endpoint"`
	Id string			`json:"id"`
}

/**
	constructor
 */
func NewQueryRegisterHandler(r *mux.Router) *QueryRegisterHandler {
	qrh := new(QueryRegisterHandler)
	qrh.router = r
	return qrh
}

/**
	register handler path
 */
func (h *QueryRegisterHandler) RegisterHandlePath() {
	h.router.Handle("/_register", h)
	logger.Debug("register path /_register")
}

func (h *QueryRegisterHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var m = QueryRegisterResponseModel{};
	var q, err = query.Get(body);

	if err == nil {
		var p, err = pipeline.CreateFromQuery(q);

		if p != nil {
			m.Result = "success";
			m.Id = p.Id;
			m.Endpoint = p.Endpoint;
		} else {
			m.Result = "fail";
			m.Message = err.Error()
		}
	} else {
		m.Result = "fail";
		m.Message = err.Error()
	}

	ret, _ := json.Marshal(m)

	rw.Write(ret)
}
