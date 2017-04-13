package handler

import (
	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
}

type HandlerInterface interface {
	RegisterHandlePath()
}

func CreateDefaultRegisteredHandlers(r *mux.Router) {
	registerHandlePath(
		NewQueryRegisterHandler(r),
		NewQueryCloseHandler(r))
}

func registerHandlePath(handlers ... HandlerInterface) {
	for _, h := range handlers {
		h.RegisterHandlePath()
	}
}