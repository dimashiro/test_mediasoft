package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/heartbeat", Heartbeat)
	return router
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
