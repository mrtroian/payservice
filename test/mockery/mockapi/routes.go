package mockapi

import (
	"fmt"
	"net/http"
)

func NewRouter(pattern string, name string) http.Handler {
	// @TODO: validate pattern
	apiEndpoint = fmt.Sprintf("/%s/", pattern)

	mux := http.NewServeMux()
	mux.Handle(apiEndpoint, http.HandlerFunc(getPayMethod))

	return mux
}
