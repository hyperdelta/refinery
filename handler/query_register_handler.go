package handler

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/hyperdelta/refinery/processor"
	"github.com/hyperdelta/refinery/config"
	"encoding/json"
	"strings"
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

	if config.Debug {
		log.Print("register path /_register, from QueryRegisterHandler")
	}
}

func (h *QueryRegisterHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var m = QueryRegisterResponseModel{};
	var q, err = processor.GetQueryObject(body);

	if err == nil {
		var _, errAgg = processor.CreateAggregatorFromQuery(q);

		if errAgg == nil {
			m.Result = "success";
		} else {
			m.Result = "fail";
			m.Message = errAgg.Error()
		}
	} else {
		m.Result = "fail";
		m.Message = err.Error()
	}

	ret, _ := json.Marshal(m)

	rw.Write(ret)

}
