package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	requestPool = sync.Pool{
		New: func() interface{} {
			return new(Request)
		},
	}
	httpClient iClient = &http.Client{}
)

type IRequest interface {
	Do() (string, error)
}

type iClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Request struct {
	url string
}

// NewRequest: return new *Request
func NewRequest(url string) *Request {
	r := requestPool.Get().(*Request)

	r.url = url

	return r
}

// (*Request).Return: return object to the object pool
// usage of (*Request) after Return invoke results in undefined behavior
func (r *Request) Return() {
	r.url = ""
	requestPool.Put(r)
}

// (*Request).Do: return response of http.Get on url
func (r *Request) Do() (string, error) {
	// resp, err := http.Get(r.url)

	req, err := http.NewRequest("GET", r.url, &bytes.Buffer{})

	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	respStruct := NewResponse()

	if err := json.Unmarshal(body, respStruct); err != nil {
		return "", err
	}

	url := respStruct.URL
	respStruct.Return()

	if len(url) <= 0 {
		return "", errors.New("url is empty")
	}

	return url, nil
}
