package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrtroian/payservice/internal/gateway"
)

var (
	apiEndpoint string
)

func ReadID(path string) (int, error) {
	if path == apiEndpoint {
		return -1, errors.New("ID not provided")
	}

	idStr := strings.TrimPrefix(path, apiEndpoint)

	id, err := strconv.Atoi(idStr)

	if err != nil || id < 1 {
		return -1, errors.New("ID is invalid")
	}

	return id, nil
}

func payMethodsHandler(w http.ResponseWriter, req *http.Request) {
	id, err := ReadID(req.URL.Path)
	resp := NewResponse()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.StatusCode = http.StatusBadRequest
		resp.Status = "400 - Bad Request"

		r, err := json.Marshal(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("api.getPayMethods: Status 500:", err)
			return
		}
		w.Write(r)
		log.Println("api.getPayMethods: Status 400:", err)
		return
	}

	gm := gateway.NewGatewayManager()
	paymentMethods, err := gm.CollectByID(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resp.StatusCode = http.StatusNotFound
		resp.Status = "404 - Not Found"

		r, err := json.Marshal(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("api.getPayMethods: Status 500:", err)
			return
		}
		w.Write(r)
		log.Println("api.getPayMethods: Status 404:", err)
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = "200 - OK"
	resp.Gateways = paymentMethods

	jsonResponse, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("api.getPayMethods: Status 500:", err)
		return
	}

	w.Write(jsonResponse)
	return
}
