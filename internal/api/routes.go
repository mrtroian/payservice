package api

import (
	"net/http"
)

func NewRouter(host, pattern string) http.Handler {
	// // @TODO: validate pattern
	// apiEndpoint = fmt.Sprintf("/%s/", pattern)
	// serverName = host

	// finalHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Println(req.URL.Path)
	// 	if strings.Compare(req.URL.Path, apiEndpoint) > 0 {
	// 		fmt.Println(req.URL.Path)
	// 		payMethodsHandler(w, req)
	// 	}
	// })
	// mux.Handle(apiEndpoint, timer(get(finalHandler)))

	return nil
}
