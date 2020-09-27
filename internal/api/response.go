package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

var responsePool = sync.Pool{
	New: func() interface{} {
		return new(Response)
	},
}

type Response struct {
	Gateways   map[string]string `json:"paymentGateways,omitempty"`
	StatusCode int               `json:"statusCode"`
	Status     string            `json:"status"`
}

// NewResponse: return new *Response
func NewResponse() *Response {
	return responsePool.Get().(*Response)
}

// (*Response).Return: return object to the object pool
// usage of (*Response) after Return invoke results in undefined behavior
func (r *Response) Return() {
	r.Gateways = nil
	r.StatusCode = 0
	r.Status = ""
	responsePool.Put(r)
}

func BadRequest() ([]byte, int) {
	resp := NewResponse()
	resp.StatusCode = http.StatusBadRequest
	resp.Status = "400 - Bad Request"

	r, err := json.Marshal(resp)

	resp.Return()

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return r, http.StatusBadRequest
}

func NotFound() ([]byte, int) {
	resp := NewResponse()
	resp.StatusCode = http.StatusNotFound
	resp.Status = "404 - Not Found"

	r, err := json.Marshal(resp)

	resp.Return()

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return r, http.StatusNotFound
}
