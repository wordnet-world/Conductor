package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter returns a *mux.Router which can be configured with
// the Routes in Routes.go and wraps all calls with the Logger
// defined in Logger.go
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		// handler = Authenticate(handler, route.Name)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
