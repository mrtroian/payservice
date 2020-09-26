package client

import "sync"

var responsePool = sync.Pool{
	New: func() interface{} {
		return new(response)
	},
}

type response struct {
	URL string `json:"paymentGatewayUrl"`
}

// NewResponse: return new *response
func NewResponse() *response {
	return responsePool.Get().(*response)
}

// (*response).Return: return object to the object pool
// usage of (*response) after Return invoke results in undefined behavior
func (r *response) Return() {
	r.URL = ""
	responsePool.Put(r)
}
