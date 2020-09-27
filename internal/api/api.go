package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mrtroian/payservice/internal/gateway"
)

type API struct {
	pattern string
	host    string
}

func NewAPI(host, pattern string) *API {
	api := new(API)

	if pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}

	api.pattern = pattern
	api.host = host

	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		resp   []byte
		status int
	)
	start := time.Now()

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
	}

	id, err := api.readID(req.URL.Path)

	if err != nil {
		status = http.StatusBadRequest
	} else {
		resp, status = api.payMethodsHandler(id)
	}

	finish := time.Since(start)
	responseTime := fmt.Sprintf("%d μs", finish.Microseconds())

	w.Header().Set("X-Server-name", api.host)
	w.Header().Add("X-Response-Time", responseTime)
	w.Header().Set("Content-Type", "application/json")

	switch status {
	case http.StatusOK:
		w.WriteHeader(status)
		w.Write(resp)
	case http.StatusNotFound:
		resp, status = NotFound()
		w.WriteHeader(status)
		w.Write(resp)
	case http.StatusBadRequest:
		resp, status = BadRequest()
		w.WriteHeader(status)
		w.Write(resp)
	}

	log.Println(req.Method, req.URL.Path, responseTime)
}

func (api *API) payMethodsHandler(id int) ([]byte, int) {
	resp := NewResponse()
	gm := gateway.NewGatewayManager()
	paymentMethods, err := gm.CollectByID(id)

	if err != nil {
		return nil, http.StatusNotFound
	}

	resp.StatusCode = http.StatusOK
	resp.Status = "200 - OK"
	resp.Gateways = paymentMethods

	jsonResponse, err := json.Marshal(resp)
	resp.Return()

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return jsonResponse, http.StatusOK
}

func (api *API) readID(path string) (int, error) {
	if path == api.pattern {
		return -1, errors.New("ID not provided")
	}

	idStr := strings.TrimPrefix(path, api.pattern)
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 1 {
		return -1, errors.New("ID is invalid")
	}

	return id, nil
}
