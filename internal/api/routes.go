package api

import (
	"fmt"
	"net/http"
)

func NewRouter(host, pattern string) http.Handler {
	// @TODO: validate pattern
	apiEndpoint = fmt.Sprintf("/%s/", pattern)
	serverName = host

	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(payMethodsHandler)
	mux.Handle(apiEndpoint, timer(get(finalHandler)))

	return mux
}
