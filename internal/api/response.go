package api

import "sync"

var responsePool = sync.Pool{
	New: func() interface{} {
		return new(Response)
	},
}

type Response struct {
	Gateways   map[string]string `json:"paymentGateways"`
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
