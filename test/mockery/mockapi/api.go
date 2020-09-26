package mockapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	apiEndpoint string
	serviceName string
)

type Response struct {
	ProductUUID string `json:"productId"`
	PayBtnURL   string `json:"paymentGatewayUrl"`
	ServiceName string `json:"serviceName"`
}

func readID(path string) (string, error) {
	if path == apiEndpoint {
		return "", errors.New("ID not provided")
	}

	id := strings.TrimPrefix(path, apiEndpoint)

	return id, nil
}

func getPayMethod(w http.ResponseWriter, req *http.Request) {
	id, err := readID(req.URL.Path)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("400 - Bad request!"))
		// @TODO: Write error
		log.Println("api.getPayMethods: Status 400:", err)
		return
	}

	r := Response{
		ProductUUID: id,
		PayBtnURL:   "https://mocked.paymentgateway.com/paywith" + serviceName,
		ServiceName: serviceName,
	}

	jsonPayMethods, err := json.Marshal(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// @TODO: Write error
		log.Println("api.getPayMethods: Status 500:", err)
		return
	}

	w.Write(jsonPayMethods)
	return
}
